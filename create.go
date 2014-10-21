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
	cs := make([]string, 0)

	for _, column := range stmt.table.Columns() {
		// Get the create syntax for each type
		ct, err := column.typ.Create(d)
		if err != nil {
			return "", err
		}
		cs = append(cs, fmt.Sprintf(`"%s" %s`, column.Name(), ct))
	}

	if stmt.table.pk != nil {
		pkc, err := stmt.table.pk.Create(d)
		if err != nil {
			return "", nil
		}
		cs = append(cs, pkc)
	}

	for _, constraint := range stmt.table.uniques {
		compiledConstraint, err := constraint.Create(d)
		if err != nil {
			return "", nil
		}
		cs = append(cs, compiledConstraint)
	}

	t := fmt.Sprintf(
		"CREATE TABLE \"%s\" (\n  %s\n);",
		stmt.table.Name,
		strings.Join(cs, ",\n  "),
	)
	return t, nil
}
