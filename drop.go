package aspect

import (
	"fmt"
)

// DropStmt is the internal representation of an DROP TABLE statement.
type DropStmt struct {
	table    *TableElem
	ifExists bool
}

// IfExists adds the IF EXISTS modifier to a DROP TABLE statement.
func (stmt DropStmt) IfExists() DropStmt {
	stmt.ifExists = true
	return stmt
}

// String outputs the parameter-less CREATE TABLE statement in a neutral
// dialect.
func (stmt DropStmt) String() string {
	c, _ := stmt.Compile(&defaultDialect{}, Params())
	return c
}

// Compile outputs the DROP TABLE statement using the given dialect and
// parameters. An error may be returned because of a pre-existing error or
// because an error occurred during compilation.
func (stmt DropStmt) Compile(d Dialect, p *Parameters) (string, error) {
	if stmt.ifExists {
		return fmt.Sprintf(`DROP TABLE IF EXISTS "%s"`, stmt.table.Name), nil
	}
	return fmt.Sprintf(`DROP TABLE "%s"`, stmt.table.Name), nil
}
