package aspect

import "fmt"

// JoinOnStmt implements a variety of joins
type JoinOnStmt struct {
	ArrayClause
	method string
	table  *TableElem
}

// String returns a default string representation of the JoinOnStmt
func (j JoinOnStmt) String() string {
	compiled, _ := j.Compile(&defaultDialect{}, Params())
	return compiled
}

// Compile compiles a JoinOnStmt
func (j JoinOnStmt) Compile(d Dialect, params *Parameters) (string, error) {
	// Compile the clauses of the join statement
	clauses, err := j.ArrayClause.Compile(d, params)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(
		` %s %s ON %s`,
		j.method,
		j.table.Compile(d, params),
		clauses,
	), nil
}
