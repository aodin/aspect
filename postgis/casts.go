package postgis

import (
	"github.com/aodin/aspect"
)

func AsGeography(c aspect.ColumnElem) aspect.ColumnElem {
	return c.SetInner(aspect.BinaryClause{
		Pre: c.Inner(),
		Sep: "::geography",
	})
}

func AsGeometry(c aspect.ColumnElem) aspect.ColumnElem {
	return c.SetInner(aspect.BinaryClause{
		Pre: c.Inner(),
		Sep: "::geometry",
	})
}
