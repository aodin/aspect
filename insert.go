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
	Stmt
	table   *TableElem
	columns []ColumnElem // TODO custom type for setter / getter operations
	args    []interface{}
	fields  fields
}

// String outputs the parameter-less INSERT statement in a neutral dialect.
func (stmt InsertStmt) String() string {
	compiled, _ := stmt.Compile(&defaultDialect{}, Params())
	return compiled
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
	if err := stmt.Error(); err != nil {
		return "", err
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
		stmt.table.Name(),
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
func (stmt *InsertStmt) trimFields(elem reflect.Value) {
	// TODO this function could be skipped if it was known that the given
	// struct has no omitempty fields
	validFields := fields{}
	for _, field := range stmt.fields {
		if !field.Exists() {
			continue
		}
		if field.HasOption(OmitEmpty) {
			var fieldElem reflect.Value = elem
			for _, index := range field.index {
				fieldElem = fieldElem.Field(index)
			}
			if isEmptyValue(fieldElem) {
				// Remove the column
				stmt.columns = removeColumn(stmt.columns, field.column)
				continue
			}
		}
		// Keep the field
		validFields = append(validFields, field)
	}
	stmt.fields = validFields
}

func (stmt *InsertStmt) argsFromValues(values Values) {
	for _, column := range stmt.columns {
		stmt.args = append(stmt.args, values[column.Name()])
	}
}

func (stmt *InsertStmt) argsFromElem(elem reflect.Value) {
	for _, field := range stmt.fields {
		var fieldElem reflect.Value = elem
		for _, index := range field.index {
			fieldElem = fieldElem.Field(index)
		}
		stmt.args = append(stmt.args, fieldElem.Interface())
	}
}

// updateColumns removes any columns that weren't matched by the fields.
func (stmt *InsertStmt) updateColumns() {
	// TODO keep actual target columns separate from requested in case
	// the statement is updated?
	for _, column := range stmt.columns {
		if !stmt.fields.HasColumn(column.Name()) {
			// TODO pass an index to prevent further nesting of iteration
			stmt.columns = removeColumn(stmt.columns, column.Name())
		}
	}
}

// Values adds parameters to the INSERT statement. If the given values do not
// match the statement's current columns, the columns will be updated.
// Valid values include structs, Values maps, or slices of structs or Values.
func (stmt InsertStmt) Values(arg interface{}) InsertStmt {
	// For now, inserts can be performed on pointers or values
	// NOTE: If auto-updating fields are required, they will need pointers
	elem := reflect.Indirect(reflect.ValueOf(arg))

	switch elem.Kind() {
	case reflect.Struct:
		// Inspect the fields of the given struct
		unaligned := SelectFieldsFromElem(elem.Type())

		// TODO function to return names of columns
		columns := make([]string, len(stmt.columns))
		for i, column := range stmt.columns {
			columns[i] = column.Name()
		}

		stmt.fields = AlignColumns(columns, unaligned)

		// If no fields were found and the number of fields matches the
		// columns requested, then insert the struct's values as is.
		if stmt.fields.Empty() && len(unaligned) == len(stmt.columns) {
			stmt.fields = unaligned
			stmt.argsFromElem(elem)
			return stmt
		}

		// Remove unmatched columns and empty values from fields
		// with the 'omitempty' option
		stmt.updateColumns()
		stmt.trimFields(elem)

		// If no fields remain after trimming, abort
		if len(stmt.fields) == 0 {
			stmt.SetError("aspect: could not match fields for INSERT - are the 'db' tags correct?")
			return stmt
		}

		// Collect the parameters
		stmt.argsFromElem(elem)

	case reflect.Slice:
		if elem.Len() == 0 {
			stmt.SetError("aspect: args cannot be set by empty slices")
			return stmt
		}
		// Slices of structs or Values are acceptable
		// TODO check kind of elem directly?
		elem0 := elem.Index(0)
		if elem0.Kind() == reflect.Struct {
			unaligned := SelectFieldsFromElem(elem.Type().Elem())

			// TODO function to return names of columns
			columns := make([]string, len(stmt.columns))
			for i, column := range stmt.columns {
				columns[i] = column.Name()
			}
			stmt.fields = AlignColumns(columns, unaligned)

			// If no fields were found and the number of fields matches the
			// columns requested, then insert the struct's values as is.
			if stmt.fields.Empty() && len(unaligned) == len(stmt.columns) {
				stmt.fields = unaligned
				for i := 0; i < elem.Len(); i++ {
					stmt.argsFromElem(elem.Index(i))
				}
				return stmt
			}

			// Remove unmatched columns and empty values from fields
			// with the 'omitempty' option
			stmt.updateColumns()
			stmt.trimFields(elem0)

			// If no fields remain after trimming, abort
			if len(stmt.fields) == 0 {
				stmt.SetError("aspect: could not match fields for INSERT - are the 'db' tags correct?")
				return stmt
			}

			// Add the parameters for each element
			for i := 0; i < elem.Len(); i++ {
				stmt.argsFromElem(elem.Index(i))
			}

			return stmt
		}

		valuesSlice, ok := arg.([]Values)
		if ok {
			if len(valuesSlice) == 0 {
				stmt.SetError(
					"aspect: cannot insert []Values of length zero",
				)
				return stmt
			}

			// Set the table columns according to the values
			var err error
			if stmt.fields, err = valuesMap(stmt, valuesSlice[0]); err != nil {
				stmt.SetError(err.Error())
				return stmt
			}

			// TODO Column names should match in each values element!
			stmt.updateColumns()

			// Add the args in the values
			for _, v := range valuesSlice {
				stmt.argsFromValues(v)
			}

			return stmt
		}

		stmt.SetError(
			"aspect: unsupported type %T for INSERT %s - values must be of type struct, Values, or a slice of either",
			arg, stmt,
		)

	case reflect.Map:
		// The only allowed map type is Values
		values, ok := arg.(Values)
		if !ok {
			stmt.SetError(
				"aspect: inserted maps must be of type Values",
			)
			return stmt
		}

		// Set the table columns according to the values
		var err error
		if stmt.fields, err = valuesMap(stmt, values); err != nil {
			stmt.SetError(err.Error())
			return stmt
		}

		// Remove unmatched columns and add args from the values
		stmt.updateColumns()
		stmt.argsFromValues(values)

	}
	return stmt
}

// TODO better way to pass columns than by using the whole statement?
// TODO if this is better generalized then it can be used with UPDATE and
// DELETE statements.
// TODO use a column set - that's all it needs - maybe move to fields?
func valuesMap(stmt InsertStmt, values Values) (fields, error) {
	fields := make(fields, len(values))
	var i int
	for column, _ := range values {
		if !stmt.HasColumn(column) {
			return nil, fmt.Errorf(
				"aspect: cannot INSERT a value with column '%s' as it has no corresponding column in the INSERT statement",
				column,
			)
		}
		fields[i] = field{column: column} // TODO set index?
	}
	return fields, nil
}

// Insert creates an INSERT statement for the given columns. There must be at
// least one column and all columns must belong to the same table.
func Insert(selection Selectable, selections ...Selectable) (stmt InsertStmt) {
	columns := make([]ColumnElem, 0)
	for _, s := range append([]Selectable{selection}, selections...) {
		if s == nil {
			stmt.SetError("aspect: insert received a nil selectable - do the columns or tables you selected exist?")
			return
		}
		columns = append(columns, s.Selectable()...)
	}

	if len(columns) < 1 {
		stmt.SetError("aspect: no columns were selected for INSERT")
		return
	}

	// The table is set from the first column
	column := columns[0]
	if column.table == nil {
		stmt.SetError(
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
			stmt.SetError("aspect: cannot INSERT to a column that does not exist - are selected columns named correctly?")
			return
		}
		if c.table != stmt.table {
			stmt.SetError("aspect: columns of an INSERT must all belong to the same table")
			return
		}
		stmt.columns = append(stmt.columns, c)
	}
	return
}
