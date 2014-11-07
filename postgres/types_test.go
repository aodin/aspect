package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aodin/aspect"
)

func TestSerial(t *testing.T) {
	assert := assert.New(t)
	expect := aspect.NewTester(t, &PostGres{})

	expect.Create("SERIAL", Serial{})
	expect.Create("SERIAL PRIMARY KEY", Serial{PrimaryKey: true})
	expect.Create(
		"SERIAL PRIMARY KEY NOT NULL",
		Serial{PrimaryKey: true, NotNull: true},
	)

	value, err := Serial{}.Validate(123)
	assert.Nil(err)
	assert.Equal(123, value)

	value, err = Serial{}.Validate("123")
	assert.Nil(err)
	assert.Equal(123, value)

	_, err = Serial{}.Validate("HEY")
	assert.NotNil(err)
}

func TestInet(t *testing.T) {
	expect := aspect.NewTester(t, &PostGres{})

	expect.Create("INET", Inet{})
}
