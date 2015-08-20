package postgres

import "github.com/aodin/aspect"

func IsEmpty(c aspect.Clause) aspect.FuncClause {
	return aspect.FuncClause{Inner: c, F: "ISEMPTY"}
}

func StringAgg(c aspect.ColumnElem, separator string) aspect.ColumnElem {
	return c.SetInner(
		aspect.FuncClause{
			Inner: aspect.ArrayClause{
				Clauses: []aspect.Clause{
					c,
					aspect.StringClause{
						Name: separator,
					},
				},
				Sep: separator,
			},
			F: "string_agg",
		},
	)
}
