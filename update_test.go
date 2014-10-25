package aspect

import (
	"testing"
)

func TestUpdate(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	expect.SQL(
		`UPDATE "users" SET "name" = $1`,
		Update(users, Values{"name": "client"}),
		"client",
	)

	values := Values{
		"name":     "admin",
		"password": "blank",
	}

	expect.SQL(
		`UPDATE "users" SET "name" = $1 AND "password" = $2 WHERE "users"."id" = $3`,
		Update(users, values).Where(users.C["id"].Equals(1)),
		"admin",
		"blank",
		1,
	)

	// The statement should have an error if the values map is empty
	stmt := Update(users, Values{})
	_, err := stmt.Compile(&defaultDialect{}, Params())
	if err == nil {
		t.Fatalf("No error returned from column-less UPDATE")
	}

	// Attempt to update values with keys that do not correspond to columns
	stmt = Update(users, Values{"nope": "what"})
	_, err = stmt.Compile(&defaultDialect{}, Params())
	if err == nil {
		t.Fatalf("no error returned from UPDATE without corresponding column")
	}
}
