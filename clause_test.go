package aspect

import (
	"testing"
)

func TestClauses(t *testing.T) {
	id := users.C["id"]

	// Binary clause
	expectedSQL(
		t,
		id.Equals(2),
		`"users"."id" = $1`,
	)

	// Array clause of binary clauses
	expectedSQL(
		t,
		AllOf(id.LessThan(5), id.GreaterThan(1)),
		`"users"."id" < $1 AND "users"."id" > $2`,
	)

	// Composite clauses
	expectedSQL(
		t,
		id.Between(2, 5),
		`"users"."id" >= $1 AND "users"."id" <= $2`,
	)
	expectedSQL(
		t,
		id.NotBetween(2, 5),
		`"users"."id" < $1 OR "users"."id" > $2`,
	)
}
