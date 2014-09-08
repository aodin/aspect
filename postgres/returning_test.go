package postgres

import (
	"github.com/aodin/aspect"
	"testing"
)

// A short test for testing that an SQL statement was compiled as expected
func expectedSQL(t *testing.T, stmt aspect.Compiles, expected string, p int) {
	params := aspect.Params()
	compiled, err := stmt.Compile(&PostGres{}, params)
	if err != nil {
		t.Error(err)
	}
	if compiled != expected {
		t.Errorf("Unexpected SQL: %s != %s", compiled, expected)
	}
	if params.Len() != p {
		t.Errorf(
			"Unexpected number of parameters for %s: %d != %d",
			expected,
			params.Len(),
			p,
		)
	}
}

var users = aspect.Table("users",
	aspect.Column("id", aspect.Integer{NotNull: true}),
	aspect.Column("name", aspect.String{Length: 32, NotNull: true}),
	aspect.Column("password", aspect.String{Length: 128}),
	aspect.PrimaryKey("id"),
)

func TestInsert(t *testing.T) {
	stmt := Insert(
		users.C["name"],
		users.C["password"],
	).Returning(
		users.C["id"],
	)

	// By default, an INSERT without values will assume a single entry
	// TODO This statement should have zero parameters
	expectedSQL(
		t,
		stmt,
		`INSERT INTO "users" ("name", "password") VALUES ($1, $2) RETURNING "users"."id"`,
		2,
	)

	// Statements with a returning clause should be unaffected
	stmt2 := Insert(
		users.C["name"],
		users.C["password"],
	)
	expectedSQL(
		t,
		stmt2,
		`INSERT INTO "users" ("name", "password") VALUES ($1, $2)`,
		2,
	)
}
