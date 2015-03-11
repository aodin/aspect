package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"

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

	// Insert multiple values
	inserts := []user{admin, {Name: "client", Password: "1234"}}
	expect.SQL(
		`INSERT INTO "users" ("name", "password") VALUES ($1, $2), ($3, $4) RETURNING "users"."id"`,
		stmt.Values(inserts),
		"admin",
		"secret",
		"client",
		"1234",
	)

	// Statements with a returning clause should be unaffected
	expect.SQL(
		`INSERT INTO "users" ("name", "password") VALUES ($1, $2)`,
		Insert(users.C["name"], users.C["password"]),
		nil,
		nil,
	)

	// Insert with a omitted field - empty and non-empty
	omitAdmin := omitID{Name: "admin", Password: "1234"}

	expect.SQL(
		`INSERT INTO "users" ("name", "password") VALUES ($1, $2)`,
		Insert(users).Values(omitAdmin),
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

	// Return all user table columns
	expect.SQL(
		`INSERT INTO "users" ("name", "password") VALUES ($1, $2) RETURNING "users"."id", "users"."name", "users"."password", "users"."is_active", "users"."created_at"`,
		Insert(users).Values(omitAdmin).Returning(users),
		"admin",
		"1234",
	)

	// Selecting a column or table that is not part of the insert table
	// should produce an error
	expect.Error(Insert(users).Values(omitAdmin).Returning(hasUUIDs))
	expect.Error(Insert(users).Values(omitAdmin).Returning(hasUUIDs.C["uuid"]))
}

func TestReturning(t *testing.T) {
	assert := assert.New(t)
	db, tx := testSetup(t)
	defer db.Close()
	defer tx.Rollback()

	tx.MustExecute(users.Create())

	clients := []user{
		{Name: "client", Password: "1234", IsActive: false},
		{Name: "member", Password: "secret", IsActive: true},
	}
	stmt := Insert(
		users.C["name"],
		users.C["password"],
		users.C["is_active"],
	).Returning(
		users.C["id"],
		users.C["created_at"],
	)
	assert.Nil(tx.QueryAll(stmt.Values(clients), &clients))

	// The IDs should be anything but zero
	assert.Equal(2, len(clients))

	assert.NotEqual(0, clients[0].ID)
	assert.Equal("client", clients[0].Name)
	assert.False(clients[0].CreatedAt.IsZero())

	assert.NotEqual(0, clients[1].ID)
	assert.Equal("member", clients[1].Name)
	assert.False(clients[1].CreatedAt.IsZero())

	// Test UUID creation
	tx.MustExecute(hasUUIDs.Create())

	u := hasUUID{Name: "what"}
	uuidStmt := Insert(hasUUIDs).Values(u).Returning(hasUUIDs)
	assert.Nil(tx.QueryOne(uuidStmt, &u))
	assert.NotEqual("", u.UUID, "UUID should have been set")
}
