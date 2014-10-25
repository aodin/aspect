package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aodin/aspect"
)

// Perform some selections against an actual postgres database
// Note: A db.json file must be set in this package

type incompleteUser struct {
	Name     string `db:"name"`
	Password string `db:"password"`
	IsActive bool   `db:"is_active"`
}

func TestSelect(t *testing.T) {
	assert := assert.New(t)

	// Connect to the database specified in the test db.json config
	conf, err := aspect.ParseConfig("./db.json")
	if err != nil {
		t.Fatalf(
			"postgres: failed to parse test configuration, test aborted: %s",
			err,
		)
	}

	db, err := aspect.Connect(conf.Driver, conf.Credentials())
	require.Nil(t, err)
	defer db.Close()

	// Perform all tests in a transaction and always rollback
	tx, err := db.Begin()
	require.Nil(t, err)
	defer tx.Rollback()

	_, err = tx.Execute(users.Create())
	require.Nil(t, err)

	// Insert users as values or pointers
	admin := user{Name: "admin", IsActive: true}
	_, err = tx.Execute(users.Insert(admin))
	require.Nil(t, err)

	var u user
	require.Nil(t, tx.QueryOne(users.Select(), &u))
	assert.Equal("admin", u.Name)
	assert.Equal(true, u.IsActive)
}
