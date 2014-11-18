package aspect

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Removing the tag will prevent the database from attempting to insert it
type managedUser struct {
	ID       int64
	Name     string `db:"name"`
	Contacts []string
	manager  string
}

type username struct {
	Name string
}

type mismatch struct {
	ID   int64  `db:"idx"`
	Name string `db:"namex"`
}

type omitID struct {
	ID       int64  `db:"id,omitempty"`
	Name     string `db:"name"`
	Password string `db:"password"`
}

func TestInsert(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	// By default, an INSERT without values will assume a single entry
	expect.SQL(
		`INSERT INTO "users" ("id", "name", "password") VALUES ($1, $2, $3)`,
		users.Insert(),
		nil,
		nil,
		nil,
	)

	stmt := Insert(users.C["name"], users.C["password"])
	expect.SQL(
		`INSERT INTO "users" ("name", "password") VALUES ($1, $2)`,
		stmt,
		nil,
		nil,
	)

	// Adding multiple values will generate a bulk insert statement
	// Structs do not need to be complete if fields are named
	admin := user{Name: "admin", Password: "secret"}
	client := user{Name: "client", Password: "1234"}
	expect.SQL(
		`INSERT INTO "users" ("name", "password") VALUES ($1, $2)`,
		stmt.Values(admin),
		"admin",
		"secret",
	)
	expect.SQL(
		`INSERT INTO "users" ("name", "password") VALUES ($1, $2), ($3, $4)`,
		stmt.Values([]user{admin, client}),
		"admin",
		"secret",
		"client",
		"1234",
	)

	// Insert with a omitted field - empty and non-empty
	expect.SQL(
		`INSERT INTO "users" ("name", "password") VALUES ($1, $2)`,
		users.Insert().Values(omitID{Name: "admin", Password: "1234"}),
		"admin",
		"1234",
	)

	expect.SQL(
		`INSERT INTO "users" ("id", "name", "password") VALUES ($1, $2, $3)`,
		users.Insert().Values(omitID{ID: 1, Name: "admin", Password: "1234"}),
		1,
		"admin",
		"1234",
	)

	// Omit should also work for multiple struct inserts
	omits := []omitID{
		omitID{Name: "admin", Password: "1234"},
		omitID{Name: "client", Password: "1234"},
	}
	expect.SQL(
		`INSERT INTO "users" ("name", "password") VALUES ($1, $2), ($3, $4)`,
		users.Insert().Values(omits),
		"admin",
		"1234",
		"client",
		"1234",
	)

	expect.SQL(
		`INSERT INTO "users" ("id", "name", "password") VALUES ($1, $2, $3)`,
		users.Insert().Values([]omitID{{ID: 1, Name: "admin", Password: "1"}}),
		1,
		"admin",
		"1",
	)

	// Complete table inserts will build dynamically from the given values
	// If a struct does not have fields that match the given columns,
	// those columns will not be included in the compiled statement
	expect.SQL(
		`INSERT INTO "users" ("id", "name", "password") VALUES ($1, $2, $3)`,
		users.Insert().Values(user{ID: 1, Name: "client"}),
		1,
		"client",
		"",
	)
	expect.SQL(
		`INSERT INTO "users" ("name") VALUES ($1)`,
		users.Insert().Values(managedUser{Name: "client"}),
		"client",
	)

	// Multiple managed users
	expect.SQL(
		`INSERT INTO "users" ("name") VALUES ($1), ($2)`,
		users.Insert().Values([]managedUser{{Name: "alice"}, {Name: "bob"}}),
		"alice",
		"bob",
	)

	// Insert a struct without db tags.
	// TODO As long as the number of exported fields matches columns
	// the statement should be allowed. As of now, no unexported fields are
	// allowed - this must be checked by using another test package
	expect.SQL(
		`INSERT INTO "users" ("name") VALUES ($1)`,
		Insert(users.C["name"]).Values(username{Name: "Solo"}),
		"Solo",
	)

	// Insert a struct without matching tags
	expect.Error(users.Insert().Values(mismatch{}))

	// Attempt to insert columns that do not exist
	expect.Error(Insert(ColumnElem{name: "what"}))
	expect.Error(Insert(users.C["id"], users.C["what"]))
	expect.Error(Insert(users.C["what"]))

	// The statement columns should be modified according to the given Values
	expect.SQL(
		`INSERT INTO "users" ("name") VALUES ($1)`,
		users.Insert().Values(Values{"name": "Hotspur"}),
		"Hotspur",
	)

	// A slice of Values is valid
	vs := []Values{
		{"name": "Totti"},
		{"name": "De Rossi"},
	}
	expect.SQL(
		`INSERT INTO "users" ("name") VALUES ($1), ($2)`,
		users.Insert().Values(vs),
		"Totti",
		"De Rossi",
	)

	// TODO what about mismatched []Values? Should having different keys
	// cause an error?

	// Extra values should cause an error
	v := Values{
		"name": "Tottenham",
		"what": "Field?",
	}
	expect.Error(users.Insert().Values(v))

	// Insert instances other than structs, values, or slices of either
	expect.Error(users.Insert().Values([]int64{1, 2, 3}))
}

func TestIsEmptyValue(t *testing.T) {
	assert := assert.New(t)
	// Expected empty values

	assert.True(isEmptyValue(reflect.ValueOf(0)))
	assert.True(isEmptyValue(reflect.ValueOf("")))
	assert.True(isEmptyValue(reflect.ValueOf(false)))
	assert.True(isEmptyValue(reflect.ValueOf(0.0)))
	assert.True(isEmptyValue(reflect.ValueOf(time.Time{})))

	assert.False(isEmptyValue(reflect.ValueOf(1)))
	assert.False(isEmptyValue(reflect.ValueOf("h")))
	assert.False(isEmptyValue(reflect.ValueOf(true)))
	assert.False(isEmptyValue(reflect.ValueOf(0.1)))
	assert.False(isEmptyValue(reflect.ValueOf(time.Now())))
}

func TestRemoveColumn(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(
		[]ColumnElem{users.C["name"], users.C["password"]},
		removeColumn(users.Columns(), "id"),
	)
	assert.Equal(
		[]ColumnElem{users.C["id"], users.C["password"]},
		removeColumn(users.Columns(), "name"),
	)
	assert.Equal(
		[]ColumnElem{users.C["id"], users.C["name"]},
		removeColumn(users.Columns(), "password"),
	)
}
