package aspect

import (
	"fmt"
)

type DropStmt struct {
	table *TableElem
}

func (stmt DropStmt) String() string {
	c, _ := stmt.Compile(&defaultDialect{}, Params())
	return c
}

// Compile creates the `DROP TABLE` statement for the given dialect.
// TODO should all statements end with semicolons?
func (stmt DropStmt) Compile(d Dialect, p *Parameters) (string, error) {
	return fmt.Sprintf(`DROP TABLE "%s"`, stmt.table.Name), nil
}
