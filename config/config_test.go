package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabase(t *testing.T) {
	c, err := ParseConfig("./example.db.json")
	assert.Nil(t, err)
	assert.Equal(t, "postgres", c.Driver)
	assert.Equal(t,
		`host=localhost port=5432 dbname=aspect_test user=postgres sslmode=disable`,
		c.Credentials(),
	)
}
