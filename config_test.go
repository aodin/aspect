package aspect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var travisCI = DatabaseConfig{
	Driver:  "postgres",
	Host:    "localhost",
	Port:    5432,
	Name:    "travis_ci_test",
	User:    "postgres",
	SSLMode: "disable",
}

// TODO Parse a json file

func TestDatabaseConfig(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("postgres", travisCI.Driver)
	assert.Equal(
		`host=localhost port=5432 dbname=travis_ci_test user=postgres sslmode=disable`,
		travisCI.Credentials(),
	)
}
