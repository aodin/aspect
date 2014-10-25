package postgres

import (
	"testing"

	"github.com/aodin/aspect"
)

type user struct {
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	Password string `db:"password"`
}

// TODO Schemas sould live in a single file

var users = aspect.Table("users",
	aspect.Column("id", aspect.Integer{NotNull: true}),
	aspect.Column("name", aspect.String{Length: 32, NotNull: true}),
	aspect.Column("password", aspect.String{Length: 128}),
	aspect.PrimaryKey("id"),
)

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
}
