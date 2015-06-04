package dbtest

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/aodin/aspect"
	"github.com/aodin/aspect/config"
)

// With returns a database connection pool and transaction using the given
// Goose config and database name.
func With(t *testing.T, path, dbname string) (*aspect.DB, aspect.Transaction) {
	// Parse the Goose db config
	c, err := config.ParseGooseDatabase(path, dbname)
	require.Nil(t, err)
	return with(t, c)
}

// WithConfig returns a database connection pool and transaction using the
// given config.Database path. It will fallback to a Travis CI spec if no
// file is found.
func WithConfig(t *testing.T, path string) (*aspect.DB, aspect.Transaction) {
	c, err := config.ParseTestConfig(path)
	require.Nil(t, err)
	return with(t, c)
}

func with(t *testing.T, c config.Database) (*aspect.DB, aspect.Transaction) {
	// Connect to the database
	conn, err := aspect.Connect(c.Driver, c.Credentials())
	require.Nil(t, err)

	// Start a transaction
	tx, err := conn.Begin()
	require.Nil(t, err)
	return conn, tx
}
