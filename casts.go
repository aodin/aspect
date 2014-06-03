package aspect

import (
	"time"
)

func (c ColumnElem) InLocation(loc *time.Location) ColumnElem {
	c.inner = BinaryClause{
		Pre:  c.inner,
		Post: &Parameter{loc.String()},
		Sep:  "::TIMESTAMP WITH TIME ZONE AT TIME ZONE ",
	}
	return c
}
