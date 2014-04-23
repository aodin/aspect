package aspect

import (
	"fmt"
)

type DeleteStmt struct {
	table *TableElem
	args  []interface{}
	cond  Clause
}

func (stmt DeleteStmt) String() string {
	compiled, _ := stmt.Compile(&PostGres{}, Params())
	return compiled
}

func (stmt DeleteStmt) Compile(d Dialect, params *Parameters) (string, error) {
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

func (stmt DeleteStmt) Where(cond Clause) DeleteStmt {
	stmt.cond = cond
	return stmt
}

func Delete(table *TableElem, args ...interface{}) DeleteStmt {
	stmt := DeleteStmt{table: table}
	// If the table has a primary key, create a where statement using
	// its columns and the values from the given args

	// TODO Bulk delete

	return stmt
}
