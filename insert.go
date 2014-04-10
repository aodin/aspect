package aspect

import (
	"fmt"
	"reflect"
	"strings"
)

type InsertStatement struct {
	table   *TableStruct
	columns []*ColumnStruct
	args    []interface{}
}

func (stmt *InsertStatement) String() string {
	return stmt.Compile()
}

func (stmt *InsertStatement) Compile() string {
	c := len(stmt.columns)
	// No columns? no statement!
	if c == 0 {
		return ""
	}

	columns := make([]string, len(stmt.columns))
	for i, column := range stmt.columns {
		columns[i] = fmt.Sprintf(`"%s"`, column.Name())
	}

	// TODO column length should divide args without remainder
	g := len(stmt.args) / c
	// If there are no arguments, default to one group
	if g == 0 {
		g = 1
	}
	parameters := make([]string, g)

	var param int
	for i := 0; i < g; i += 1 {
		group := make([]string, c)
		for j := 0; j < c; j += 1 {
			param += 1
			// TODO Parameters are dialect specific
			group[j] = fmt.Sprintf("$%d", param)
		}
		parameters[i] = fmt.Sprintf(`(%s)`, strings.Join(group, ", "))
	}

	// TODO Bulk insert syntax is dialect specific
	return fmt.Sprintf(
		`INSERT INTO "%s" (%s) VALUES %s`,
		stmt.table.Name,
		strings.Join(columns, ", "),
		strings.Join(parameters, ", "),
	)
}

func (stmt *InsertStatement) Args() []interface{} {
	return stmt.args
}

func (stmt *InsertStatement) Execute() (string, error) {
	return stmt.Compile(), nil
}

// There must be at least one arg
func (stmt *InsertStatement) Values(arg interface{}, args ...interface{}) *InsertStatement {

	var l int // Expected length of each argument, set by first arg
	elem := reflect.Indirect(reflect.ValueOf(arg))

	// TODO What if there are existing values attached to the stmt?

	// TODO Allow slice types
	if elem.Kind() == reflect.Struct {
		// TODO They must match the types and length of any previous args
		l = elem.NumField()
		for i := 0; i < l; i += 1 {
			stmt.args = append(stmt.args, elem.Field(i).Interface())
		}
		for _, arg := range args {
			e := reflect.Indirect(reflect.ValueOf(arg))
			for i := 0; i < l; i += 1 {
				stmt.args = append(stmt.args, e.Field(i).Interface())
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

func InsertTableValues(t *TableStruct, arg interface{}, args ...interface{}) *InsertStatement {
	stmt := &InsertStatement{
		table:   t,
		columns: t.Columns(),
		args:    make([]interface{}, 0),
	}
	return stmt.Values(arg, args...)
}
