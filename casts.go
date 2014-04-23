package aspect

import (
	"time"
)

func (c ColumnElem) InLocation(loc *time.Location) ColumnElem {
	c.inner = BinaryClause{
		pre:  c.inner,
		post: &Parameter{loc.String()},
		sep:  "::TIMESTAMP WITH TIME ZONE AT TIME ZONE ",
	}
	return c
}
