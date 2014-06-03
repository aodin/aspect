package postgis

import (
	"github.com/aodin/aspect"
)

func AsGeography(inner aspect.Clause) aspect.BinaryClause {
	return aspect.BinaryClause{
		Pre:  inner,
		Post: aspect.StringClause{"geography"},
		Sep:  "::",
	}
}
