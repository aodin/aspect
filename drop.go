package aspect

import (
	"fmt"
)

type DropStmt struct {
	table    *TableElem
	ifExists bool
}

func (stmt DropStmt) IfExists() DropStmt {
	stmt.ifExists = true
	return stmt
}

// String outputs a parameter-less and dialect neutral DROP TABLE statement.
func (stmt DropStmt) String() string {
	c, _ := stmt.Compile(&defaultDialect{}, Params())
	return c
}

// Compile creates the `DROP TABLE` statement for the given dialect.
// TODO should all statements end with semicolons?
func (stmt DropStmt) Compile(d Dialect, p *Parameters) (string, error) {
	if stmt.ifExists {
		return fmt.Sprintf(`DROP TABLE IF EXISTS "%s"`, stmt.table.Name), nil
	}
	return fmt.Sprintf(`DROP TABLE "%s"`, stmt.table.Name), nil
}
