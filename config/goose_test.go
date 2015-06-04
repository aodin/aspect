package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTestYAML(t *testing.T) {
	c, err := ParseGooseDatabase("./example.dbconf.yml", "test")
	assert.Nil(t, err)
	assert.Equal(t, "postgres", c.Driver)
	assert.Equal(t,
		"host=localhost port=5432 dbname=db_test user=test password=bad sslmode=disable",
		c.Credentials(),
	)
}
