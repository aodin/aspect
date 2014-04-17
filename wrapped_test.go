package aspect

import (
	"testing"
)

func TestAggregates(t *testing.T) {
	// COUNT()
	count := Count(views.C["id"])
	expectedPostGres(t, count, `COUNT("views"."id")`, 0)

	// DATE_PART()
	datePart := DatePart(views.C["timestamp"], "quarter")
	// TODO expect one parameter?
	expectedPostGres(
		t,
		datePart,
		`DATE_PART('quarter', "views"."timestamp")`,
		0,
	)
}
