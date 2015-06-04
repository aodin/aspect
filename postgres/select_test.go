package postgres

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aodin/aspect"
	"github.com/aodin/aspect/dbtest"
)

// Perform some selections against an actual postgres database
// Note: A db.json file must be set in this package

type incompleteUser struct {
	Name     string `db:"name"`
	Password string `db:"password"`
	IsActive bool   `db:"is_active"`
}

// Select incomplete structs
type testUser struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Password  string    `db:"password"`
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
	Contacts  []string
	manager   string
}

func TestSelect(t *testing.T) {
	// Connect to the database specified in the test db.json config
	// Default to the Travis CI settings if no file is found
	conn, tx := dbtest.WithConfig(t, "./db.json")

	// Perform all tests in a transaction and always rollback
	defer conn.Close()
	defer tx.Rollback()

	_, err := tx.Execute(users.Create())
	require.Nil(t, err)

	// Insert users as values or pointers
	admin := user{Name: "admin", IsActive: true}
	stmt := aspect.Insert(
		users.C["name"],
		users.C["password"],
		users.C["is_active"],
	)
	_, err = tx.Execute(stmt.Values(admin))
	require.Nil(t, err)

	_, err = tx.Execute(stmt.Values(&admin))
	require.Nil(t, err)

	var u user
	require.Nil(t, tx.QueryOne(users.Select(), &u))
	assert.Equal(t, "admin", u.Name)
	assert.True(t, u.IsActive)

	// Select using a returning clause
	client := user{Name: "client", Password: "1234", IsActive: false}
	returningStmt := Insert(
		users.C["name"],
		users.C["password"],
		users.C["is_active"],
	).Returning(
		users.C["id"],
		users.C["created_at"],
	)
	require.Nil(t, tx.QueryOne(returningStmt.Values(client), &client))

	// The ID and time should be anything but zero
	assert.NotEqual(t, int64(0), client.ID)
	assert.False(t, client.CreatedAt.IsZero())

	// Select into a struct that has extra columns
	// TODO Skip unexported fields
	var extraField testUser
	require.Nil(t, tx.QueryOne(users.Select().Where(users.C["name"].Equals("client")), &extraField))
	assert.Equal(t, "client", extraField.Name)
	assert.False(t, extraField.IsActive)

	// Query multiple users
	var extraFields []testUser
	assert.Nil(t, tx.QueryAll(users.Select(), &extraFields))
	assert.Equal(t, 3, len(extraFields))

	// Query ids directly
	var ids []int64
	orderBy := aspect.Select(users.C["id"]).OrderBy(users.C["id"].Desc())
	assert.Nil(t, tx.QueryAll(orderBy, &ids))
	assert.Equal(t, 3, len(ids))
	var prev int64
	for _, id := range ids {
		if prev != 0 {
			if prev < id {
				t.Errorf("id results returned out of order")
			}
		}
	}

	// Scan into a Values map - all fields should be returned
	selectByName := users.Select().OrderBy(users.C["name"])
	values := aspect.Values{}
	assert.Nil(t, tx.QueryOne(selectByName, values))
	assert.Equal(t, 5, len(values))
	assert.NotEqual(t, int64(0), values["id"])
	assert.Equal(t, []byte("admin"), values["name"]) // Yup, strings are []byte
	assert.EqualValues(t, true, values["is_active"])
	assert.Equal(t, []byte{}, values["password"])

	// Scan into a slice of Values maps
	var allValues []aspect.Values
	assert.Nil(t, tx.QueryAll(selectByName, &allValues))
	assert.Equal(t, 3, len(allValues))

	values = allValues[2]
	assert.Equal(t, 5, len(values))
	assert.NotEqual(t, int64(0), values["id"])
	assert.Equal(t, []byte("client"), values["name"])
	assert.EqualValues(t, false, values["is_active"])
	assert.Equal(t, []byte("1234"), values["password"])

	// TODO Test duplicate names in a table join

	// TODO destination types that don't match the result

	// Update
	// ------

	updateStmt := users.Update().Values(
		aspect.Values{"name": "HELLO", "password": "BYE"},
	).Where(
		users.C["id"].Equals(client.ID),
	)
	result, err := tx.Execute(updateStmt)
	require.Nil(t, err)

	rowsAffected, err := result.RowsAffected()
	assert.Nil(t, err)
	assert.EqualValues(t, 1, rowsAffected)

	// Delete
	// ------

	result, err = tx.Execute(users.Delete())
	require.Nil(t, err)

	rowsAffected, err = result.RowsAffected()
	assert.Nil(t, err)
	assert.EqualValues(t, 3, rowsAffected)
}

type u1 struct {
	ID   int64  `db:"id,omitempty"`
	Name string `db:"name"`
}

var u1s = aspect.Table("u1s",
	aspect.Column("id", Serial{PrimaryKey: true, NotNull: true}),
	aspect.Column("name", aspect.String{NotNull: true}),
)

func TestSelect_Omitempty(t *testing.T) {
	// Connect to the database specified in the test db.json config
	// Default to the Travis CI settings if no file is found
	conn, tx := dbtest.WithConfig(t, "./db.json")

	// Perform all tests in a transaction and always rollback
	defer conn.Close()
	defer tx.Rollback()

	_, err := tx.Execute(u1s.Create())
	require.Nil(t, err)

	u := u1{Name: "admin"}
	stmt := Insert(u1s).Values(u).Returning(u1s.C["id"])
	err = tx.QueryOne(stmt, &u.ID)
	require.Nil(t, err)
	assert.Equal(t, "admin", u.Name)

	// ID should be auto-assigned
	assert.True(t, u.ID > 0)
}
