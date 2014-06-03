package aspect

import (
	"testing"
)

func TestFunctions(t *testing.T) {
	// COUNT()
	count := Count(views.C["id"])
	expectedSQL(t, count, `COUNT("views"."id")`, 0)

	// DATE()
	date := DateOf(views.C["timestamp"])
	expectedSQL(t, date, `DATE("views"."timestamp")`, 0)

	// DATE_PART()
	datePart := DatePart(views.C["timestamp"], "quarter")
	expectedSQL(
		t,
		datePart,
		`DATE_PART('quarter', "views"."timestamp")`,
		0,
	)
}
