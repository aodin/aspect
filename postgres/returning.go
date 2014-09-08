package postgres

import (
	"errors"
	"fmt"
	"github.com/aodin/aspect"
	"reflect"
	"strings"
)

type InsertStmt struct {
	table     *aspect.TableElem
	columns   []aspect.ColumnElem
	returning []aspect.ColumnElem
	args      []interface{}
	err       error
	l         int
	alias     []string
}

var (
	ErrNoColumns = errors.New("postgres: statments must have associated columns")
)

func (stmt InsertStmt) String() string {
	compiled, _ := stmt.Compile(&PostGres{}, aspect.Params())
	return compiled
}

func (stmt InsertStmt) Error() error {
	return stmt.err
}

func (stmt InsertStmt) Compile(d aspect.Dialect, params *aspect.Parameters) (string, error) {
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
			p := &aspect.Parameter{stmt.args[param]}
			group[j], _ = p.Compile(d, params)
			param += 1
		}
		parameters[i] = fmt.Sprintf(`(%s)`, strings.Join(group, ", "))
	}

	compiled := fmt.Sprintf(
		`INSERT INTO "%s" (%s) VALUES %s`,
		stmt.table.Name,
		strings.Join(columns, ", "),
		strings.Join(parameters, ", "),
	)

	if len(stmt.returning) > 0 {
		compiled += fmt.Sprintf(
			" RETURNING %s",
			strings.Join(stmt.CompileColumns(d, params, stmt.returning), ", "),
		)
	}
	return compiled, nil
}

func (stmt InsertStmt) Returning(columns ...aspect.ColumnElem) InsertStmt {
	for _, column := range columns {
		if column.Table() != stmt.table {
			stmt.err = fmt.Errorf(
				"postgres: all columns must belong to the same table",
			)
			break
		}
		stmt.returning = append(stmt.returning, column)
	}
	return stmt
}

// Iterate through the struct fields and see if the tags match
// the given column names.
// Return the field names matching their respective columns.
// The field tag takes precendence over the name.
func fieldAlias(cs []aspect.ColumnElem, i interface{}) []string {
	// Get the type of the interface pointer
	t := reflect.TypeOf(i)
	if t.Kind() != reflect.Ptr {
		t = reflect.PtrTo(t)
	}
	// TODO Confirm that the given interface is a struct

	alias := make([]string, len(cs))
	// For each field, try the tag name, then the field name
	elem := t.Elem()
	for ci, column := range cs {
		name := column.Name()
		for i := 0; i < elem.NumField(); i += 1 {
			f := elem.Field(i)
			fname := f.Name
			tag := f.Tag.Get("db")
			if tag == name || name == fname {
				alias[ci] = fname
				break
			}
		}
	}
	// TODO Were all the columns matched?
	return alias
}

// Get the value of the field struct by name
func (stmt *InsertStmt) argsByName(elem reflect.Value) {
	for _, n := range stmt.alias {
		stmt.args = append(stmt.args, elem.FieldByName(n).Interface())
	}
}

// Read every value of the struct in order
func (stmt *InsertStmt) argsByIndex(elem reflect.Value) {
	for i := 0; i < stmt.l; i += 1 {
		stmt.args = append(stmt.args, elem.Field(i).Interface())
	}
}

// There must be at least one arg
func (stmt InsertStmt) Values(args interface{}) InsertStmt {
	// For now, inserts can be performed on pointers or values
	// TODO If auto-updating fields are required, they will need pointers
	elem := reflect.Indirect(reflect.ValueOf(args))

	// TODO What if there are existing values attached to the stmt?
	// Skip these checks if there is a populated alias / l

	switch elem.Kind() {
	case reflect.Struct:
		// If the number of columns does not match the number of fields,
		// attempt to build an alias object
		stmt.l = elem.NumField()
		if stmt.l != len(stmt.columns) {
			// TODO confirm that the alias is fully populated
			// TODO fieldAlias should be a method that can only be set once
			stmt.alias = fieldAlias(stmt.columns, args)
			stmt.argsByName(elem)
		} else {
			stmt.argsByIndex(elem)
		}
	case reflect.Slice:
		sliceLen := elem.Len()
		if sliceLen < 0 {
			// TODO There must be values - a delayed error?
			return stmt
		}
		firstElem := elem.Index(0)

		// TODO Slice elements must be structs for now
		stmt.l = firstElem.NumField()

		if stmt.l != len(stmt.columns) {
			// TODO confirm that the alias is fully populated
			stmt.alias = fieldAlias(stmt.columns, firstElem.Interface())

			// Add every slice element to the args by field name
			for i := 0; i < sliceLen; i++ {
				stmt.argsByName(elem.Index(i))
			}
		} else {
			// And every slice elem's fields to the args
			for i := 0; i < sliceLen; i++ {
				stmt.argsByIndex(elem.Index(i))
			}
		}
	case reflect.Map:
		// Set the table columns according to the values
		stmt.columns = make([]aspect.ColumnElem, 0)

		// Cast back to Values
		values, ok := args.(aspect.Values)
		if !ok {
			stmt.err = fmt.Errorf(
				"aspect: to insert maps they must of type aspect.Values",
			)
			return stmt
		}
		for _, key := range values.Keys() {
			column, exists := stmt.table.C[key]
			if !exists {
				stmt.err = fmt.Errorf(
					`aspect: no column "%s" exists in the table "%s"`,
					key,
					stmt.table.Name,
				)
				return stmt
			}
			stmt.columns = append(stmt.columns, column)
			stmt.args = append(stmt.args, values[key])
		}
		// TODO Allow []Values, which must have the same columns
	}
	return stmt
}

// TODO are errors required?
func (stmt InsertStmt) CompileColumns(d aspect.Dialect, params *aspect.Parameters, columns []aspect.ColumnElem) []string {
	names := make([]string, len(columns))
	for i, c := range columns {
		names[i], _ = c.Compile(d, params)
	}
	return names
}

// Insert creates an INSERT statement for the given columns. There must be at
// least one column and all columns must belong to the same table.
func Insert(c aspect.ColumnElem, columns ...aspect.ColumnElem) InsertStmt {
	stmt := InsertStmt{
		table:   c.Table(),
		columns: []aspect.ColumnElem{c},
		args:    make([]interface{}, 0),
	}

	for _, column := range columns {
		if column.Table() != stmt.table {
			// TODO How to best handle delayed errors?
			stmt.err = fmt.Errorf("All columns must belong to the same table")
			break
		}
		stmt.columns = append(stmt.columns, column)
	}
	return stmt
}
