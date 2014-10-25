package aspect

import (
	"testing"
	"time"
)

func TestClauses(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	id := users.C["id"]

	// Column clause
	expect.SQL(`"users"."id"`, ColumnClause{table: users, name: id.Name()})

	// Binary clause
	expect.SQL(`"users"."id" = $1`, id.Equals(2), 2)

	// Load the Denver timezone
	denver, err := time.LoadLocation("America/Denver")
	if err != nil {
		t.Fatal(err)
	}
	expect.SQL(
		`"views"."timestamp"::TIMESTAMP WITH TIME ZONE AT TIME ZONE $1`,
		views.C["timestamp"].InLocation(denver),
		denver.String(),
	)

	// Array clause of binary clauses
	expect.SQL(
		`"users"."id" < $1 AND "users"."id" > $2`,
		AllOf(id.LessThan(5), id.GreaterThan(1)),
		5,
		1,
	)

	// Composite clauses
	expect.SQL(
		`"users"."id" >= $1 AND "users"."id" <= $2`,
		id.Between(2, 5),
		2,
		5,
	)
	expect.SQL(
		`"users"."id" < $1 OR "users"."id" > $2`,
		id.NotBetween(2, 5),
		2,
		5,
	)
	expect.SQL(
		`"users"."id" IN ($1, $2)`,
		id.In([]int64{1, 5}),
		1,
		5,
	)

	// Unary clauses
	expect.SQL(
		`"users"."id" IS NULL`,
		id.IsNull(),
	)
	expect.SQL(
		`"users"."id" IS NOT NULL`,
		id.IsNotNull(),
	)
}
