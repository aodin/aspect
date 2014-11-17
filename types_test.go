package aspect

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestText(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	expect.Create("TEXT", Text{})
}

func TestString(t *testing.T) {
	assert := assert.New(t)
	expect := NewTester(t, &defaultDialect{})

	expect.Create("VARCHAR", String{})
	expect.Create(
		"VARCHAR(128) PRIMARY KEY NOT NULL UNIQUE",
		String{Length: 128, PrimaryKey: true, NotNull: true, Unique: true},
	)

	// Test Type methods
	value, err := String{}.Validate("HEY")
	assert.Nil(err)
	assert.Equal("HEY", value)

	_, err = String{}.Validate(123)
	assert.NotNil(err)

	_, err = String{Length: 3}.Validate("HELLO")
	assert.NotNil(err)
}

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

func TestTimestamp(t *testing.T) {
	assert := assert.New(t)
	expect := NewTester(t, &defaultDialect{})

	expect.Create(
		"TIMESTAMP DEFAULT (now() at time zone 'utc')",
		Timestamp{Default: "now() at time zone 'utc'"},
	)
	expect.Create(
		"TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (now() at time zone 'utc')",
		Timestamp{
			WithTimezone: true,
			NotNull:      true,
			Default:      "now() at time zone 'utc'",
		},
	)

	d := time.Date(2014, 1, 1, 12, 0, 0, 0, time.UTC)
	value, err := Timestamp{}.Validate(d)
	assert.Nil(err)
	assert.Equal(d, value)

	_, err = Timestamp{}.Validate(123)
	assert.NotNil(err)
}

func TestDate(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	expect.Create(
		"DATE PRIMARY KEY NOT NULL UNIQUE",
		Date{PrimaryKey: true, NotNull: true, Unique: true},
	)
}

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
