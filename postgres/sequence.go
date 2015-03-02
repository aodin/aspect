package postgres

import (
	"fmt"

	"github.com/aodin/aspect"
)

type Sequence string

type AlterSeqStmt struct {
	sequence Sequence
	clause   aspect.Clause
}

func (stmt AlterSeqStmt) Compile(d aspect.Dialect, ps *aspect.Parameters) (string, error) {
	// A clause is required
	if stmt.clause == nil {
		return "", fmt.Errorf(
			"postgres: ALTER SEQUENCE statements require a clause",
		)
	}

	// Compile the internal clause
	cc, err := stmt.clause.Compile(d, ps)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`ALTER SEQUENCE "%s" %s`, stmt.sequence, cc), nil
}

func (stmt AlterSeqStmt) RenameTo(name string) AlterSeqStmt {
	stmt.clause = aspect.BinaryClause{
		Sep:  "RENAME TO ",
		Post: aspect.StringClause{Name: name},
	}
	return stmt
}

func (stmt AlterSeqStmt) RestartWith(n int) AlterSeqStmt {
	stmt.clause = aspect.BinaryClause{
		Sep:  "RESTART WITH ",
		Post: aspect.IntClause{D: n},
	}
	return stmt
}

func (stmt AlterSeqStmt) String() string {
	output, _ := stmt.Compile(&PostGres{}, aspect.Params())
	return output
}

func AlterSequence(sequence Sequence) (stmt AlterSeqStmt) {
	stmt.sequence = sequence
	return
}
