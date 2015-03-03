package aspect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	assert := assert.New(t)
	expect := NewTester(t, &defaultDialect{})

	expect.Create("VARCHAR", String{})
	expect.Create(
		"VARCHAR(128) PRIMARY KEY NOT NULL UNIQUE",
		String{Length: 128, PrimaryKey: true, NotNull: true, Unique: true},
	)

	expect.Create("VARCHAR DEFAULT ''", String{Default: Blank})

	// Test Type methods
	value, err := String{}.Validate("HEY")
	assert.Nil(err)
	assert.Equal("HEY", value)

	_, err = String{}.Validate(123)
	assert.NotNil(err)

	_, err = String{Length: 3}.Validate("HELLO")
	assert.NotNil(err)
}
