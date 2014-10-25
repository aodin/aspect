package aspect

import (
	"testing"
)

func TestFunctions(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})
	expect.SQL(`AVG("views"."id")`, Avg(views.C["id"]))
	expect.SQL(`SUM("views"."id")`, Sum(views.C["id"]))
	expect.SQL(`COUNT("views"."id")`, Count(views.C["id"]))
	expect.SQL(`DATE("views"."timestamp")`, DateOf(views.C["timestamp"]))
	expect.SQL(`LOWER("views"."id")`, Lower(views.C["id"]))
	expect.SQL(`MAX("views"."id")`, Max(views.C["id"]))
	expect.SQL(
		`DATE_PART('quarter', "views"."timestamp")`,
		DatePart(views.C["timestamp"], "quarter"),
	)
}
