package aspect

import (
	"testing"
)

func TestFunctions(t *testing.T) {
	expectedSQL(t, Avg(views.C["id"]), `AVG("views"."id")`, 0)
	expectedSQL(t, Sum(views.C["id"]), `SUM("views"."id")`, 0)
	expectedSQL(t, Count(views.C["id"]), `COUNT("views"."id")`, 0)
	expectedSQL(
		t,
		DateOf(views.C["timestamp"]),
		`DATE("views"."timestamp")`,
		0,
	)
	expectedSQL(t, Max(views.C["id"]), `MAX("views"."id")`, 0)
	expectedSQL(
		t,
		DatePart(views.C["timestamp"], "quarter"),
		`DATE_PART('quarter', "views"."timestamp")`,
		0,
	)
}
