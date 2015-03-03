package aspect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDouble(t *testing.T) {
	assert := assert.New(t)
	expect := NewTester(t, &defaultDialect{})

	expect.Create("DOUBLE PRECISION", Double{})
	expect.Create(
		"DOUBLE PRECISION PRIMARY KEY NOT NULL UNIQUE",
		Double{PrimaryKey: true, NotNull: true, Unique: true},
	)

	value, err := Double{}.Validate(123.456)
	assert.Nil(err)
	assert.Equal(123.456, value)

	value, err = Double{}.Validate(123)
	assert.Nil(err)
	assert.Equal(float64(123), value)

	value, err = Double{}.Validate("123.456")
	assert.Nil(err)
	assert.Equal(123.456, value)

	_, err = Double{}.Validate("HEY")
	assert.NotNil(err)
}
