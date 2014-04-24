package aspect

import (
	"fmt"
)

type JoinStmt struct {
	method    string
	table     *TableElem
	pre, post ColumnElem
}

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
