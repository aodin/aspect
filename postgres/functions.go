package postgres

import "github.com/aodin/aspect"

func IsEmpty(c aspect.Clause) aspect.FuncClause {
	return aspect.FuncClause{Inner: c, F: "ISEMPTY"}
}
