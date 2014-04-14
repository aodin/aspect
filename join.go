package aspect

import (
	"fmt"
)

type JoinStmt struct {
	table *TableStruct
	pre   ColumnElement
	post  ColumnElement
}

func (j *JoinStmt) Compile(d Dialect, params *Parameters) (string, error) {
	prec, err := j.pre.Compile(d, params)
	if err != nil {
		return "", err
	}
	postc, err := j.post.Compile(d, params)
	if err != nil {
		return "", err
	}
	compiled := fmt.Sprintf(
		` JOIN "%s" ON %s = %s`,
		j.table.Name,
		prec,
		postc,
	)
	return compiled, nil
}
