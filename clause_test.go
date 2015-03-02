package aspect

import (
	"testing"
	"time"
)

func TestClauses(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	id := users.C["id"]
	name := users.C["name"]

	// Column equality - used in JOIN ... ON ... statements
	expect.SQL(
		`"users"."id" = "users"."name"`,
		id.Equals(name),
	)

	// Column clause
	column := ColumnClause{table: users, name: name.Name()}
	expect.SQL(`"users"."name"`, column)

	// String and Int clause TODO remove these as the skip parameterization
	expect.SQL(`'name'`, StringClause{Name: "name"})
	expect.SQL(`3`, IntClause{D: 3})

	// Func clause
	expect.SQL(`LOWER("users"."name")`, FuncClause{Inner: column, F: "LOWER"})

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
		`("users"."id" < $1 AND "users"."id" > $2)`,
		AllOf(id.LessThan(5), id.GreaterThan(1)),
		5,
		1,
	)
	expect.SQL(
		`(("users"."id" < $1 AND "users"."id" > $2) OR "users"."id" = $3)`,
		AnyOf(AllOf(id.LessThan(5), id.GreaterThan(1)), id.Equals(7)),
		5,
		1,
		7,
	)

	// Composite clauses
	expect.SQL(
		`("users"."id" >= $1 AND "users"."id" <= $2)`,
		id.Between(2, 5),
		2,
		5,
	)
	expect.SQL(
		`("users"."id" < $1 OR "users"."id" > $2)`,
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
