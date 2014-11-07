package aspect

import (
	"fmt"
	"strings"
)

// Creatable is the interface that clauses used in CREATE TABLE statements
// must implement.
type Creatable interface {
	Create(d Dialect) (string, error)
}

// CreateStmt is the internal representation of an CREATE TABLE statement.
type CreateStmt struct {
	table *TableElem
}

// String outputs the parameter-less CREATE TABLE statement in a neutral
// dialect.
func (stmt CreateStmt) String() string {
	c, _ := stmt.Compile(&defaultDialect{}, Params())
	return c
}

// Compile outputs the CREATE TABLE statement using the given dialect and
// parameters. An error may be returned because of a pre-existing error or
// because an error occurred during compilation.
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
