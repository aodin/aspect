package aspect

import (
	"fmt"
)

type DeleteStatement struct {
	table *TableStruct
	args  []interface{}
	cond  Clause
}

func (stmt *DeleteStatement) String() string {
	compiled, _ := stmt.Compile(&PostGres{}, Params())
	return compiled
}

func (stmt *DeleteStatement) Compile(d Dialect, params *Parameters) (string, error) {
	compiled := fmt.Sprintf(`DELETE FROM "%s"`, stmt.table.Name)

	// TODO Add any existing arguments to the parameters

	if stmt.cond != nil {
		cc, err := stmt.cond.Compile(d, params)
		if err != nil {
			return "", err
		}
		compiled += fmt.Sprintf(" WHERE %s", cc)
	}
	return compiled, nil
}

func (stmt *DeleteStatement) Where(cond Clause) *DeleteStatement {
	stmt.cond = cond
	return stmt
}

func Delete(table *TableStruct, args ...interface{}) *DeleteStatement {
	stmt := &DeleteStatement{table: table}
	// If the table has a primary key, create a where statement using
	// its columns and the values from the given args

	if len(args) < 2 {

	}
	// TODO Bulk delete

	return stmt
}
