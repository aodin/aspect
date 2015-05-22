package aspect

import "testing"

func TestStmt(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	// Add a conditional with an existing clause
	stmt := users.Select().Where(users.C["id"].Equals(1))
	expect.SQL(
		`SELECT "users"."id", "users"."name", "users"."password" FROM "users" WHERE "users"."id" = $1`,
		stmt,
		1,
	)

	stmt.AddConditional(users.C["name"].DoesNotEqual(""))
	expect.SQL(
		`SELECT "users"."id", "users"."name", "users"."password" FROM "users" WHERE ("users"."id" = $1 AND "users"."name" != $2)`,
		stmt,
		1,
		"",
	)

	// Add a conditional without an existing clause
	stmt = users.Select()
	stmt.AddConditional(users.C["name"].DoesNotEqual(""))
	expect.SQL(
		`SELECT "users"."id", "users"."name", "users"."password" FROM "users" WHERE "users"."name" != $1`,
		stmt,
		"",
	)

}
