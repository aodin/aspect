package aspect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReal(t *testing.T) {
	assert := assert.New(t)
	expect := NewTester(t, &defaultDialect{})

	expect.Create(
		"REAL PRIMARY KEY NOT NULL UNIQUE",
		Real{PrimaryKey: true, NotNull: true, Unique: true},
	)

	value, err := Real{}.Validate(123.456)
	assert.Nil(err)
	assert.Equal(123.456, value)

	value, err = Real{}.Validate(123)
	assert.Nil(err)
	assert.Equal(float64(123), value)

	value, err = Real{}.Validate("123.456")
	assert.Nil(err)
	assert.Equal(123.456, value)

	_, err = Real{}.Validate("HEY")
	assert.NotNil(err)
}
