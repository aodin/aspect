package aspect

import (
	"fmt"
	"strings"
)

type CreateStmt struct {
	table *TableElem
}

func (stmt CreateStmt) String() string {
	c, _ := stmt.Compile(&defaultDialect{}, Params())
	return c
}

func (stmt CreateStmt) Compile(d Dialect, p *Parameters) (string, error) {
	// Compiled elements
	compiled := make([]string, len(stmt.table.creates))

	var err error
	for i, create := range stmt.table.creates {
		if compiled[i], err = create.Create(d); err != nil {
			return "", err
		}
	}

	return fmt.Sprintf(
		"CREATE TABLE \"%s\" (\n  %s\n);",
		stmt.table.Name,
		strings.Join(compiled, ",\n  "),
	), nil
}
