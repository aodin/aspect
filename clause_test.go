package aspect

import (
	"testing"
	"time"
)

func TestClauses(t *testing.T) {
	id := users.C["id"]

	// Column clause
	expectedSQL(
		t,
		ColumnClause{table: users, name: id.Name()},
		`"users"."id"`,
		0,
	)

	// Binary clause
	expectedSQL(
		t,
		id.Equals(2),
		`"users"."id" = $1`,
		1,
	)

	// Load the Denver timezone
	denver, err := time.LoadLocation("America/Denver")
	if err != nil {
		t.Fatal(err)
	}
	expectedSQL(
		t,
		views.C["timestamp"].InLocation(denver),
		`"views"."timestamp"::TIMESTAMP WITH TIME ZONE AT TIME ZONE $1`,
		1,
	)

	// Array clause of binary clauses
	expectedSQL(
		t,
		AllOf(id.LessThan(5), id.GreaterThan(1)),
		`"users"."id" < $1 AND "users"."id" > $2`,
		2,
	)

	// Composite clauses
	expectedSQL(
		t,
		id.Between(2, 5),
		`"users"."id" >= $1 AND "users"."id" <= $2`,
		2,
	)
	expectedSQL(
		t,
		id.NotBetween(2, 5),
		`"users"."id" < $1 OR "users"."id" > $2`,
		2,
	)
	expectedSQL(
		t,
		id.In([]int64{1, 5}),
		`"users"."id" IN ($1, $2)`,
		2,
	)

	// Unary clauses
	expectedSQL(
		t,
		id.IsNull(),
		`"users"."id" IS NULL`,
		0,
	)
	expectedSQL(
		t,
		id.IsNotNull(),
		`"users"."id" IS NOT NULL`,
		0,
	)
}
