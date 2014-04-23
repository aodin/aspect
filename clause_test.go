package aspect

import (
	"testing"
)

func TestClauses(t *testing.T) {
	id := users.C["id"]

	// Column clause
	expectedPostGres(
		t,
		ColumnClause{table: users, name: id.Name()},
		`"users"."id"`,
		0,
	)

	// Binary clause
	expectedPostGres(
		t,
		id.Equals(2),
		`"users"."id" = $1`,
		1,
	)

	// Array clause of binary clauses
	expectedPostGres(
		t,
		AllOf(id.LessThan(5), id.GreaterThan(1)),
		`"users"."id" < $1 AND "users"."id" > $2`,
		2,
	)

	// Composite clauses
	expectedPostGres(
		t,
		id.Between(2, 5),
		`"users"."id" >= $1 AND "users"."id" <= $2`,
		2,
	)
	expectedPostGres(
		t,
		id.NotBetween(2, 5),
		`"users"."id" < $1 OR "users"."id" > $2`,
		2,
	)
}
