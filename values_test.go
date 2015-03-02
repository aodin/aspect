package aspect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValues(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})
	expect.SQL(
		`"age" = $1, "name" = $2`,
		Values{"age": 19, "name": "Chad"},
		19,
		"Chad",
	)
}

func TestValues_Keys(t *testing.T) {
	assert := assert.New(t)
	values := Values{"name": "Chad", "age": 19}
	assert.Equal([]string{"age", "name"}, values.Keys())
}

func TestValues_Diff(t *testing.T) {
	assert := assert.New(t)

	a := Values{"1": 1, "2": 2, "3": 3}
	b := Values{"2": 3, "3": 3}
	assert.Equal(Values{"1": 1, "2": 2}, a.Diff(b))
}
