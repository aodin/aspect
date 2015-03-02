package aspect

import (
	"testing"
)

func TestUpdate(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	// Values do not need to be attached to produce an UPDATE statement. It
	// will default to all columns in the table with nil parameters.
	expect.SQL(
		`UPDATE "users" SET "id" = $1, "name" = $2, "password" = $3`,
		users.Update(),
		nil,
		nil,
		nil,
	)

	expect.SQL(
		`UPDATE "users" SET "name" = $1`,
		users.Update().Values(Values{"name": "client"}),
		"client",
	)

	values := Values{
		"name":     "admin",
		"password": "blank",
	}

	// With Where
	expect.SQL(
		`UPDATE "users" SET "name" = $1, "password" = $2 WHERE "users"."id" = $3`,
		Update(users).Values(values).Where(users.C["id"].Equals(1)),
		"admin",
		"blank",
		1,
	)

	expect.SQL(
		`UPDATE "users" SET "name" = $1 WHERE ("users"."id" = $2 AND "users"."name" = $3)`,
		Update(users).Values(Values{"name": "client"}).Where(
			users.C["id"].Equals(1),
			users.C["name"].Equals("admin"),
		),
		"client",
		1,
		"admin",
	)

	// The statement should have an error if the values map is empty
	expect.Error(users.Update().Values(Values{}))

	// Attempt to update values with keys that do not correspond to columns
	expect.Error(Update(users).Values(Values{"nope": "what"}))
}
