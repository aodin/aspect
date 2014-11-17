package postgres

import (
	"testing"

	"github.com/aodin/aspect"
)

type omitID struct {
	ID       int64  `db:"id,omitempty"`
	Name     string `db:"name"`
	Password string `db:"password"`
}

func TestInsert(t *testing.T) {
	expect := aspect.NewTester(t, &PostGres{})

	stmt := Insert(
		users.C["name"],
		users.C["password"],
	).Returning(
		users.C["id"],
	)

	// By default, an INSERT without values will assume a single entry
	expect.SQL(
		`INSERT INTO "users" ("name", "password") VALUES ($1, $2) RETURNING "users"."id"`,
		stmt,
		nil,
		nil,
	)

	// Adding values should set parameters
	admin := user{Name: "admin", Password: "secret"}
	expect.SQL(
		`INSERT INTO "users" ("name", "password") VALUES ($1, $2) RETURNING "users"."id"`,
		stmt.Values(admin),
		"admin",
		"secret",
	)

	// Statements with a returning clause should be unaffected
	expect.SQL(
		`INSERT INTO "users" ("name", "password") VALUES ($1, $2)`,
		Insert(users.C["name"], users.C["password"]),
		nil,
		nil,
	)

	// Insert with a omitted field - empty and non-empty
	expect.SQL(
		`INSERT INTO "users" ("name", "password") VALUES ($1, $2)`,
		Insert(users).Values(omitID{Name: "admin", Password: "1234"}),
		"admin",
		"1234",
	)

	expect.SQL(
		`INSERT INTO "users" ("id", "name", "password") VALUES ($1, $2, $3)`,
		Insert(users).Values(omitID{ID: 1, Name: "admin", Password: "1234"}),
		1,
		"admin",
		"1234",
	)
}
