package aspect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseConfig(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("postgres", travisCI.Driver)
	assert.Equal(
		`host=localhost port=5432 dbname=travis_ci_test user=postgres sslmode=disable`,
		travisCI.Credentials(),
	)
}
