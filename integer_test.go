package aspect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInteger(t *testing.T) {
	assert := assert.New(t)
	expect := NewTester(t, &defaultDialect{})

	expect.Create("INTEGER", Integer{})
	expect.Create(
		"INTEGER PRIMARY KEY NOT NULL UNIQUE",
		Integer{PrimaryKey: true, NotNull: true, Unique: true},
	)

	assert.Equal(false, Integer{}.IsPrimaryKey())
	assert.Equal(false, Integer{}.IsUnique())

	assert.Equal(true, Integer{PrimaryKey: true}.IsPrimaryKey())
	assert.Equal(true, Integer{Unique: true}.IsUnique())

	value, err := Integer{}.Validate(123)
	assert.Nil(err)
	assert.Equal(123, value)

	value, err = Integer{}.Validate(123.000)
	assert.Nil(err)
	assert.Equal(123, value)

	value, err = Integer{}.Validate("123")
	assert.Nil(err)
	assert.Equal(123, value)

	_, err = Integer{}.Validate(123.456)
	assert.NotNil(err)

	_, err = Integer{}.Validate("HEY")
	assert.NotNil(err)
}
