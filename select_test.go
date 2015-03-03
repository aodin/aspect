package aspect

import (
	"testing"
)

func TestSelect(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	// All three of these select statements should produce the same output
	selects := []Compiles{
		users.Select(),
		Select(users),
		Select(users.C["id"], users.C["name"], users.C["password"]),
	}
	for _, s := range selects {
		expect.SQL(
			`SELECT "users"."id", "users"."name", "users"."password" FROM "users"`,
			s,
		)
	}

	// Only select a subset of the available columns
	stmt := Select(users.C["id"], users.C["name"])
	expect.SQL(
		`SELECT "users"."id", "users"."name" FROM "users"`,
		stmt,
	)

	// Add an ORDER BY
	expect.SQL(
		`SELECT "users"."id", "users"."name" FROM "users" ORDER BY "users"."id" DESC`,
		stmt.OrderBy(users.C["id"].Desc()),
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

	// Mutiple conditionals will joined with AND by default
	expect.SQL(
		`SELECT "users"."name" FROM "users" WHERE ("users"."id" = $1 AND "users"."name" = $2)`,
		Select(users.C["name"]).Where(
			users.C["id"].Equals(1),
			users.C["name"].Equals("admin"),
		),
		1,
		"admin",
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

	// Select a column that doesn't exist
	expect.Error(Select(users.C["what"]))
}

func TestSelectTable(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	expect.SQL(
		`SELECT "users"."id", "users"."name", "users"."password", "views"."timestamp" FROM "users" JOIN "views" ON "views"."user_id" = "users"."id"`,
		users.Select(views.C["timestamp"]).JoinOn(
			views, views.C["user_id"].Equals(users.C["id"]),
		),
	)
}
