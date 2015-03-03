package aspect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoolean(t *testing.T) {
	assert := assert.New(t)
	expect := NewTester(t, &defaultDialect{})

	expect.Create("BOOLEAN", Boolean{})
	expect.Create("BOOLEAN NOT NULL", Boolean{NotNull: true})
	expect.Create(
		"BOOLEAN NOT NULL DEFAULT FALSE",
		Boolean{NotNull: true, Default: False},
	)

	value, err := Boolean{}.Validate(true)
	assert.Nil(err)
	assert.Equal(true, value)

	_, err = Boolean{}.Validate(123)
	assert.NotNil(err)
}
