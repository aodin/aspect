package aspect

import (
	"testing"
)

func TestFieldAlias(t *testing.T) {
	columns := []ColumnElem{users.C["name"], users.C["password"]}

	// Get the alias fields for the users struct
	var u user
	alias := fieldAlias(columns, &u)
	if len(alias) != len(columns) {
		t.Fatalf("Expected alias of length %d, received %d", len(columns), len(alias))
	}
	if alias[0] != "Name" {
		t.Errorf("Unexpected alias: %s != Name", alias[0])
	}
	if alias[1] != "Password" {
		t.Errorf("Unexpected alias: %s != Password", alias[1])
	}

	// Alias should work with addresses or values
	alias = fieldAlias(columns, u)
	if len(alias) != len(columns) {
		t.Fatalf("Expected alias of length %d, received %d", len(columns), len(alias))
	}
	if alias[0] != "Name" {
		t.Errorf("Unexpected alias: %s != Name", alias[0])
	}
	if alias[1] != "Password" {
		t.Errorf("Unexpected alias: %s != Password", alias[1])
	}
}

func TestInsert(t *testing.T) {
	stmt := Insert(users.C["name"], users.C["password"])

	// By default, an INSERT without values will assume a single entry
	// TODO This statement should have zero parameters
	expectedSQL(
		t,
		stmt,
		`INSERT INTO "users" ("name", "password") VALUES ($1, $2)`,
		2,
	)

	// Adding multiple values will generate a bulk insert statement
	// Structs do not need to be complete if fields are named
	admin := user{Name: "admin", Password: "secret"}
	client := user{Name: "client", Password: "1234"}

	single := stmt.Values(admin)
	expectedSQL(
		t,
		single,
		`INSERT INTO "users" ("name", "password") VALUES ($1, $2)`,
		2,
	)

	bulk := stmt.Values([]user{admin, client})
	expectedSQL(
		t,
		bulk,
		`INSERT INTO "users" ("name", "password") VALUES ($1, $2), ($3, $4)`,
		4,
	)
}
