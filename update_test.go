package aspect

import (
	"testing"
)

func TestUpdate(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	expect.SQL(
		`UPDATE "users" SET "name" = $1`,
		users.Update().Values(Values{"name": "client"}),
		"client",
	)

	values := Values{
		"name":     "admin",
		"password": "blank",
	}

	expect.SQL(
		`UPDATE "users" SET "name" = $1 AND "password" = $2 WHERE "users"."id" = $3`,
		Update(users).Values(values).Where(users.C["id"].Equals(1)),
		"admin",
		"blank",
		1,
	)

	// The statement should have an error if the values map is empty
	expect.Error(users.Update().Values(Values{}))

	// Attempt to update values with keys that do not correspond to columns
	expect.Error(Update(users).Values(Values{"nope": "what"}))
}
