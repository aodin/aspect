package postgres

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/aodin/aspect"
)

// Return a database connection pool and transaction
func testSetup(t *testing.T) (*aspect.DB, aspect.Transaction) {
	// Connect to the database specified in the test db.json config
	// Default to the Travis CI settings if no file is found
	conf, err := aspect.ParseTestConfig("./db.json")
	require.Nil(t, err)

	db, err := aspect.Connect(conf.Driver, conf.Credentials())
	require.Nil(t, err)

	tx, err := db.Begin()
	require.Nil(t, err)
	return db, tx
}

func TestDB(t *testing.T) {
	// Connect to the database specified in the test db.json config
	// Default to the Travis CI settings if no file is found
	conf, err := aspect.ParseTestConfig("./db.json")
	if err != nil {
		t.Fatalf(
			"postgres: failed to parse test configuration, test aborted: %s",
			err,
		)
	}

	db, err := aspect.Connect(conf.Driver, conf.Credentials())
	require.Nil(t, err)
	defer db.Close()

	_, err = db.Execute(users.Drop().IfExists())
	require.Nil(t, err)

	// Perform test twice
	for i := 0; i < 2; i++ {
		tx, err := db.Begin()
		require.Nil(t, err)

		// Start a fake transaction that will implement the connection passed
		// to resources / controllers
		fakeTX := aspect.FakeTx(tx)

		innerTX, err := fakeTX.Begin()
		require.Nil(t, err)

		_, err = innerTX.Execute(users.Create())
		require.Nil(t, err)
		innerTX.Commit()

		// Rollback the real transaction
		tx.Rollback()
	}
}
