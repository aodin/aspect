package aspect

import (
	"testing"
)

func TestUpdate(t *testing.T) {
	stmt := Update(users, Values{"name": "client"})
	expectedSQL(
		t,
		stmt,
		`UPDATE "users" SET "name" = $1`,
		1,
	)

	values := Values{
		"name":     "admin",
		"password": "blank",
	}

	stmt = Update(users, values).Where(users.C["id"].Equals(1))
	expectedSQL(
		t,
		stmt,
		`UPDATE "users" SET "name" = $1 AND "password" = $2 WHERE "users"."id" = $3`,
		3,
	)

	// The statement should have an error if the values map is empty
	stmt = Update(users, Values{})
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
