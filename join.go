package aspect

import (
	"fmt"
)

type JoinOnStmt struct {
	ArrayClause
	method string
	table  *TableElem
}

func (j JoinOnStmt) String() string {
	compiled, _ := j.Compile(&defaultDialect{}, Params())
	return compiled
}

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

// JoinStmt is an internal representation of a JOIN.
// It is broken and deprecated.
type JoinStmt struct {
	method    string
	table     *TableElem
	pre, post ColumnElem
}

// Compile will compile the JOIN statement according to the given dialect.
func (j JoinStmt) Compile(d Dialect, params *Parameters) (string, error) {
	prec, err := j.pre.Compile(d, params)
	if err != nil {
		return "", err
	}
	postc, err := j.post.Compile(d, params)
	if err != nil {
		return "", err
	}
	compiled := fmt.Sprintf(
		` %s "%s" ON %s = %s`,
		j.method,
		j.table.Name,
		prec,
		postc,
	)
	return compiled, nil
}
