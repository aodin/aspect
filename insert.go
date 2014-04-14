package aspect

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type InsertStatement struct {
	table   *TableElem
	columns []*ColumnStruct
	args    []interface{}
	err     error
}

var (
	ErrNoColumns = errors.New("aspect: an INSERT must have associated columns")
)

func (stmt *InsertStatement) String() string {
	compiled, _ := stmt.Compile(&PostGres{}, Params())
	return compiled
}

func (stmt *InsertStatement) Compile(d Dialect, params *Parameters) (string, error) {
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
		return "", fmt.Errorf("Size mismatch between arguments and columns: %d is not a multiple of %d", len(stmt.args), c)
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

// Iterate through the struct fields and see if the tags match
// the given column names.
// Return the field names matching their respective columns.
// The field tag takes precendence over the name.
func fieldAlias(cs []*ColumnStruct, i interface{}) []string {
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

// There must be at least one arg
func (stmt *InsertStatement) Values(arg interface{}, args ...interface{}) *InsertStatement {

	var l int // Expected length of each argument, set by first arg
	elem := reflect.Indirect(reflect.ValueOf(arg))

	// TODO What if there are existing values attached to the stmt?

	// TODO Allow slice types
	if elem.Kind() == reflect.Struct {
		// TODO Args must match the types and length of any previous args
		l = elem.NumField()

		// If the number of columns does not match the number of fields,
		// attempt to build an alias object
		if l != len(stmt.columns) {
			// TODO confirm that the alias is fully populated
			alias := fieldAlias(stmt.columns, arg)

			// Get the value of the field struct by name
			for _, n := range alias {
				stmt.args = append(stmt.args, elem.FieldByName(n).Interface())
			}
			for _, arg := range args {
				e := reflect.Indirect(reflect.ValueOf(arg))
				for _, n := range alias {
					stmt.args = append(stmt.args, e.FieldByName(n).Interface())
				}
			}
		} else {
			// Read every value of the struct in order
			for i := 0; i < l; i += 1 {
				stmt.args = append(stmt.args, elem.Field(i).Interface())
			}
			for _, arg := range args {
				e := reflect.Indirect(reflect.ValueOf(arg))
				for i := 0; i < l; i += 1 {
					stmt.args = append(stmt.args, e.Field(i).Interface())
				}
			}
		}
	} else {
		// Single value arguments. How many columns are specified by the
		// insert statement? Group the arguments by the number of columns.
	}
	return stmt
}

// There must be at least one column
func Insert(column *ColumnStruct, columns ...*ColumnStruct) *InsertStatement {
	// All columns must belong to the same table
	stmt := &InsertStatement{
		table:   column.table,
		columns: []*ColumnStruct{column},
		args:    make([]interface{}, 0),
	}

	for _, column := range columns {
		if column.table != stmt.table {
			// TODO How to handle delayed errors?
			continue
		}
		stmt.columns = append(stmt.columns, column)
	}
	return stmt
}

func InsertTableValues(t *TableElem, arg interface{}, args ...interface{}) *InsertStatement {
	stmt := &InsertStatement{
		table:   t,
		columns: t.Columns(),
		args:    make([]interface{}, 0),
	}
	return stmt.Values(arg, args...)
}
