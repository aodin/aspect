package aspect

import (
	"testing"
)

var views = Table("views",
	Column("id", Integer{}),
	Column("user_id", Integer{}),
	Column("timestamp", Timestamp{}),
	PrimaryKey("id"),
)

func TestSelect(t *testing.T) {
	// All three of these select statements should produce the same output
	s := Select(users.C["id"], users.C["name"])
	expectedSQL(t, s, `SELECT "users"."id", "users"."name" FROM "users"`)

	// Add an ORDER BY
	expectedSQL(
		t,
		s.OrderBy(users.C["id"].Desc()),
		`SELECT "users"."id", "users"."name" FROM "users" ORDER BY "users"."id" DESC`,
	)

	// Build a GROUP BY statement with sorting using an aggregate
	expectedSQL(
		t,
		Select(views.C["user_id"], Count(views.C["timestamp"])).GroupBy(views.C["user_id"]).OrderBy(Count(views.C["timestamp"]).Desc()),
		`SELECT "views"."user_id", COUNT("views"."timestamp") FROM "views" GROUP BY "views"."user_id" ORDER BY COUNT("views"."timestamp") DESC`,
	)

	// Add a conditional
	expectedSQL(
		t,
		Select(users.C["name"]).Where(users.C["id"].Equals(1)),
		`SELECT "users"."name" FROM "users" WHERE "users"."id" = $1`,
	)
}
