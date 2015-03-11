package aspect

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

var (
	ErrNoColumns = errors.New(
		"aspect: attempt to create a statement with zero columns",
	)
)

// InsertStmt is the internal representation of an INSERT statement.
type InsertStmt struct {
	table   *TableElem
	columns []ColumnElem // TODO custom type for setter / getter operations
	args    []interface{}
	err     error // TODO common error handling struct
	alias   map[string]Field
}

// String outputs the parameter-less INSERT statement in a neutral dialect.
func (stmt InsertStmt) String() string {
	compiled, _ := stmt.Compile(&defaultDialect{}, Params())
	return compiled
}

// Error returns any error attached to this statement
func (stmt InsertStmt) Error() error {
	return stmt.err
}

// SetError attaches an error to this statment
func (stmt *InsertStmt) SetError(err error) {
	stmt.err = err
}

// Table returns the table of this statement
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

// Compile outputs the INSERT statement using the given dialect and parameters.
// An error may be returned because of a pre-existing error or because
// an error occurred during compilation.
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
			stmt.args[i] = nil
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

// isEmptyValue is from Go's encoding/json package: encode.go
// Copyright 2010 The Go Authors. All rights reserved.
// TODO what about pointer fields?
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.String:
		return v.String() == ""
	case reflect.Bool:
		return !v.Bool()
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	case reflect.Struct:
		t, ok := v.Interface().(time.Time)
		if ok {
			return t.IsZero()
		}
	}
	return false
}

// TODO error if there was no match?
func removeColumn(columns []ColumnElem, name string) []ColumnElem {
	for i, column := range columns {
		if column.name == name {
			return append(columns[:i], columns[i+1:]...)
		}
	}
	return columns
}

// A field marked omitempty can cause the removal of a column, only
// to have another value not have an empty value for that field
func (stmt *InsertStmt) removeEmptyColumns(elem reflect.Value) {
	// TODO this function could be skipped if it was known that the given
	// struct has no omitempty fields
	for name, field := range stmt.alias {
		if field.OmitEmpty && isEmptyValue(elem.FieldByName(field.Name)) {
			// Remove the column
			stmt.columns = removeColumn(stmt.columns, name)
			continue
		}
	}
}

// Since alias is a map, columns must be read in order and then aliased
func (stmt *InsertStmt) argsByAlias(elem reflect.Value) {
	for _, column := range stmt.columns {
		alias := stmt.alias[column.Name()].Name
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
		stmt.args = append(stmt.args, values[stmt.alias[column.Name()].Name])
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

		// Remove empty values
		s.removeEmptyColumns(elem)

		// Add the args using the created field map
		// Drop args with empty values if they have the option "omitempty"
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

			// Remove empty values using the first elem
			s.removeEmptyColumns(elem0)

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

		s.err = fmt.Errorf(
			"aspect: unsupported type %T for INSERT %s - values must be of type struct, Values, or a slice of either",
			arg,
			s,
		)

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
func Insert(selection Selectable, selections ...Selectable) (stmt InsertStmt) {
	columns := make([]ColumnElem, 0)
	for _, s := range append([]Selectable{selection}, selections...) {
		if s == nil {
			stmt.err = fmt.Errorf("aspect: insert received a nil selectable - do the columns or tables you selected exist?")
			return
		}
		columns = append(columns, s.Selectable()...)
	}

	if len(columns) < 1 {
		stmt.err = fmt.Errorf(
			"aspect: no columns were selected for INSERT",
		)
		return
	}

	// The table is set from the first column
	column := columns[0]
	if column.table == nil {
		stmt.err = fmt.Errorf(
			"aspect: attempting to INSERT to a column unattached to a table",
		)
		return
	}
	stmt.table = column.table

	// Prepend the first column
	for _, c := range columns {
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
