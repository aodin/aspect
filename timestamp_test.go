package aspect

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
