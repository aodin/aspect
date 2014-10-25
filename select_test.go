package aspect

import (
	"testing"
)

func TestSelect(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	// All three of these select statements should produce the same output
	s := Select(users.C["id"], users.C["name"])
	expect.SQL(
		`SELECT "users"."id", "users"."name" FROM "users"`,
		s,
	)

	// Add an ORDER BY
	expect.SQL(
		`SELECT "users"."id", "users"."name" FROM "users" ORDER BY "users"."id" DESC`,
		s.OrderBy(users.C["id"].Desc()),
	)

	// Build a GROUP BY statement with sorting using an aggregate
	expect.SQL(
		`SELECT "views"."user_id", COUNT("views"."timestamp") FROM "views" GROUP BY "views"."user_id" ORDER BY COUNT("views"."timestamp") DESC`,
		Select(views.C["user_id"], Count(views.C["timestamp"])).GroupBy(views.C["user_id"]).OrderBy(Count(views.C["timestamp"]).Desc()),
	)

	// Add a conditional
	expect.SQL(
		`SELECT "users"."name" FROM "users" WHERE "users"."id" = $1`,
		Select(users.C["name"]).Where(users.C["id"].Equals(1)),
		1,
	)

	// Test limit
	expect.SQL(
		`SELECT "users"."name" FROM "users" LIMIT 1`,
		Select(users.C["name"]).Limit(1),
	)

	// Test Offset
	expect.SQL(
		`SELECT "users"."name" FROM "users" OFFSET 1`,
		Select(users.C["name"]).Offset(1),
	)
}
