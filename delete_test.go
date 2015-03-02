package aspect

import (
	"testing"
)

func TestDelete(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	// Test a complete delete
	expect.SQL(
		`DELETE FROM "users"`,
		users.Delete(),
	)

	// Test a delete with a WHERE
	expect.SQL(
		`DELETE FROM "users" WHERE "users"."id" = $1`,
		users.Delete().Where(users.C["id"].Equals(1)),
		1,
	)

	expect.SQL(
		`DELETE FROM "users" WHERE ("users"."id" = $1 AND "users"."name" = $2)`,
		users.Delete().Where(
			users.C["id"].Equals(1),
			users.C["name"].Equals("admin"),
		),
		1,
		"admin",
	)

	// Delete by a schema's declared primary key
	admin := user{ID: 1, Name: "admin", Password: "secret"}
	client := user{ID: 2, Name: "client", Password: "secret"}
	expect.SQL(
		`DELETE FROM "users" WHERE "users"."id" = $1`,
		users.Delete().Values(admin),
		1,
	)

	expect.SQL(
		`DELETE FROM "users" WHERE "users"."id" IN ($1, $2)`,
		users.Delete().Values([]user{admin, client}),
		1,
		2,
	)

	// Attempt to delete an empty slice
	expect.Error(users.Delete().Values([]user{}))

	// Attempt to delete a value with no pk (it has to be "id")
	var what = struct {
		ID int64
	}{
		ID: 0,
	}
	expect.Error(users.Delete().Values(what))

	// No pk specified
	expect.Error(singleColumn.Delete().Values(what))

	// Composite pk
	expect.Error(edges.Delete().Values(what))

	// Non-struct slice
	expect.Error(users.Delete().Values([]int64{1, 2, 3}))
}
