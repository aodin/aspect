package aspect

import (
	"fmt"
)

type DropStmt struct {
	table *TableElem
}

func (stmt DropStmt) String() string {
	c, _ := stmt.Compile(&PostGres{}, Params())
	return c
}

func (stmt DropStmt) Compile(d Dialect, p *Parameters) (string, error) {
	return fmt.Sprintf(`DROP TABLE "%s"`, stmt.table.Name), nil
}
