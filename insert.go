package aspect

import (
	"fmt"
	"reflect"
	"strings"
)

type InsertStatement struct {
	table   *TableStruct
	columns []*ColumnStruct
	args    [][]interface{}
}

func (stmt *InsertStatement) String() string {
	return stmt.Compile()
}

func (stmt *InsertStatement) Compile() string {
	columns := make([]string, len(stmt.columns))
	for i, column := range stmt.columns {
		columns[i] = fmt.Sprintf(`"%s"`, column.Name)
	}

	parameters := make([]string, len(stmt.args))
	var c int
	// TODO Or just len(stmt.columns) * stmt.args
	for i, group := range stmt.args {
		args := make([]string, len(group))
		for j, _ := range group {
			// TODO Parameters are also dialect specific
			c += 1
			args[j] = fmt.Sprintf("$%d", c)
		}
		parameters[i] = fmt.Sprintf(`(%s)`, strings.Join(args, ", "))
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
	args := make([]interface{}, 0)
	for _, group := range stmt.args {
		for _, arg := range group {
			args = append(args, arg)
		}
	}
	return args
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
		inter := make([]interface{}, l)

		for i, _ := range inter {
			inter[i] = elem.Field(i).Interface()
		}
		stmt.args = [][]interface{}{inter}

		for _, arg := range args {
			e := reflect.Indirect(reflect.ValueOf(arg))
			inter := make([]interface{}, l)
			for i, _ := range inter {
				inter[i] = e.Field(i).Interface()
			}
			stmt.args = append(stmt.args, inter)
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
		args:    make([][]interface{}, 0),
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
		columns: t.Selectable(),
		args:    make([][]interface{}, len(args)),
	}
	return stmt.Values(arg, args...)
}
