package aspect

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	ErrNoColumns = errors.New(
		"aspect: attempt to create a statement with zero columns",
	)
)

type InsertStmt struct {
	table   *TableElem
	columns []ColumnElem
	args    []interface{}
	err     error
	alias   map[string]string
}

func (stmt InsertStmt) String() string {
	compiled, _ := stmt.Compile(&defaultDialect{}, Params())
	return compiled
}

func (stmt InsertStmt) Error() error {
	return stmt.err
}

func (stmt InsertStmt) SetError(err error) {
	stmt.err = err
}

func (stmt InsertStmt) Table() *TableElem {
	return stmt.table
}

func (stmt *InsertStmt) AppendColumn(c ColumnElem) {
	stmt.columns = append(stmt.columns, c)
}

func (stmt InsertStmt) HasColumn(name string) bool {
	for _, column := range stmt.columns {
		if column.Name() == name {
			return true
		}
	}
	return false
}

func (stmt InsertStmt) Compile(d Dialect, params *Parameters) (string, error) {
	// Check for delayed errors
	if stmt.err != nil {
		return "", stmt.err
	}

	c := len(stmt.columns)
	// No columns? no statement!
	if c == 0 {
		return "", ErrNoColumns
	}

	columns := make([]string, len(stmt.columns))
	for i, column := range stmt.columns {
		columns[i] = fmt.Sprintf(`"%s"`, column.Name())
	}

	// Column length must divide args without remainder
	if len(stmt.args)%c != 0 {
		return "", fmt.Errorf(
			`aspect: size mismatch between arguments and columns: %d is not a multiple of %d`,
			len(stmt.args),
			c,
		)
	}

	g := len(stmt.args) / c
	// If there are no arguments, default to one group and create
	// placeholder values
	if g == 0 {
		g = 1
		stmt.args = make([]interface{}, c)
		for i, _ := range stmt.args {
			stmt.args[i] = struct{}{}
		}
	}
	parameters := make([]string, g)

	var param int
	for i := 0; i < g; i += 1 {
		group := make([]string, c)
		for j := 0; j < c; j += 1 {
			// Parameters are dialect specific
			// TODO errors
			p := &Parameter{stmt.args[param]}
			group[j], _ = p.Compile(d, params)
			param += 1
		}
		parameters[i] = fmt.Sprintf(`(%s)`, strings.Join(group, ", "))
	}

	// TODO Bulk insert syntax is dialect specific
	return fmt.Sprintf(
		`INSERT INTO "%s" (%s) VALUES %s`,
		stmt.table.Name,
		strings.Join(columns, ", "),
		strings.Join(parameters, ", "),
	), nil
}

func (stmt *InsertStmt) argsByAlias(elem reflect.Value) {
	// Since alias is a map, columns must be read in order and then aliased
	for _, column := range stmt.columns {
		alias := stmt.alias[column.Name()]
		stmt.args = append(stmt.args, elem.FieldByName(alias).Interface())
	}
}

// Read every value of the struct in order
func (stmt *InsertStmt) argsByIndex(elem reflect.Value) {
	for i, _ := range stmt.columns {
		stmt.args = append(stmt.args, elem.Field(i).Interface())
	}
}

func (stmt *InsertStmt) argsByValues(values Values) {
	// Since alias is a map, columns must be read in order and then aliased
	for _, column := range stmt.columns {
		stmt.args = append(stmt.args, values[stmt.alias[column.Name()]])
	}
}

// setColumns removes any columns that weren't matched by the alias.
func (stmt *InsertStmt) setColumns() {
	// TODO keep actual target columns separate from requested in case
	// the statement is updated?
	if len(stmt.alias) != len(stmt.columns) {
		var matched []ColumnElem
		for _, column := range stmt.columns {
			if _, exists := stmt.alias[column.Name()]; exists {
				matched = append(matched, column)
			}
		}
		stmt.columns = matched
	}
}

// Values adds parameters to the INSERT statement. If the given values do not
// match the statement's current columns, the columns will be updated.
// Valid values include structs, Values maps, or slices of structs or Values.
func (s InsertStmt) Values(arg interface{}) InsertStmt {
	// TODO reset args everytime this method is called?
	// For now, inserts can be performed on pointers or values
	// NOTE: If auto-updating fields are required, they will need pointers
	elem := reflect.Indirect(reflect.ValueOf(arg))

	switch elem.Kind() {
	case reflect.Struct:
		// Map the column names to fields on the given struct
		if s.alias, s.err = fieldMap(s.columns, arg); s.err != nil {
			return s
		}
		// If no columns were detected and the number of fields matches the
		// columns requested, then insert the struct's values as is.
		if len(s.alias) == 0 && len(s.columns) == elem.NumField() {
			s.argsByIndex(elem)
			return s
		} else if len(s.alias) == 0 {
			s.err = fmt.Errorf("aspect: cannot insert given struct, no fields match columns - were `db` struct tags set?")
			return s
		}

		// Remove unmatched columns
		s.setColumns()

		// Add the args using the created field map
		s.argsByAlias(elem)

	case reflect.Slice:
		if elem.Len() < 1 {
			s.err = fmt.Errorf("aspect: args cannot be set by empty slices")
			return s
		}
		// Slices of structs or Values are acceptable
		// TODO check kind of elem directly?
		elem0 := elem.Index(0)
		if elem0.Kind() == reflect.Struct {
			// TODO Remove code duplication
			if s.alias, s.err = fieldMap(s.columns, elem0.Interface()); s.err != nil {
				return s
			}
			// If no columns were detected and the number of fields matches the
			// columns requested, then insert the struct's values as is.
			// TODO This does not ignore unexported fields - it should
			if len(s.alias) == 0 && len(s.columns) == elem.NumField() {
				for i := 0; i < elem.Len(); i++ {
					s.argsByIndex(elem.Index(i))
				}
				return s
			} else if len(s.alias) == 0 {
				s.err = fmt.Errorf("aspect: cannot insert given struct, no fields match columns - were `db` struct tags set?")
				return s
			}

			// Remove unmatched columns
			s.setColumns()

			// Add the args using the created field map
			for i := 0; i < elem.Len(); i++ {
				s.argsByAlias(elem.Index(i))
			}

			return s
		}

		valuesSlice, ok := arg.([]Values)
		if ok {
			if len(valuesSlice) < 1 {
				s.err = fmt.Errorf(
					"aspect: cannot insert []Values of length zero",
				)
				return s
			}

			// Set the table columns according to the first values
			if s.alias, s.err = valuesMap(s, valuesSlice[0]); s.err != nil {
				return s
			}

			// Remove unmatched columns
			s.setColumns()

			// Add the args in the values
			// TODO what to do about varying values in the slice?
			for _, v := range valuesSlice {
				s.argsByValues(v)
			}

			return s
		}

		s.err = fmt.Errorf("aspect: unsupported type for INSERT %s - values must be of type struct, Values, or a slice of either")

	case reflect.Map:
		// The only allowed map type is Values
		values, ok := arg.(Values)
		if !ok {
			s.err = fmt.Errorf(
				"aspect: inserted maps must be of type Values",
			)
			return s
		}

		// Set the table columns according to the values
		if s.alias, s.err = valuesMap(s, values); s.err != nil {
			return s
		}

		// Remove unmatched columns
		s.setColumns()

		// Add the args in the values
		s.argsByValues(values)

	}
	return s
}

// Insert creates an INSERT statement for the given columns. There must be at
// least one column and all columns must belong to the same table.
func Insert(column ColumnElem, columns ...ColumnElem) (stmt InsertStmt) {
	// The table is set from the first column
	if column.table == nil {
		stmt.err = fmt.Errorf("aspect: attempting to INSERT to a column unattached to a table")
		return
	}
	stmt.table = column.table

	// Prepend the first column
	for _, c := range append([]ColumnElem{column}, columns...) {
		// Columns must have a name or they wouldn't exist (and probably
		// don't if the name is missing - the most common error case will
		// be a mistyped name in table.C)
		if c.Name() == "" {
			stmt.err = fmt.Errorf("aspect: cannot INSERT to a column that does not exist - are selected columns named correctly?")
			return
		}
		if c.table != stmt.table {
			stmt.err = fmt.Errorf("aspect: columns of an INSERT must all belong to the same table")
			return
		}
		stmt.columns = append(stmt.columns, c)
	}
	return
}
