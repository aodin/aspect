package sqlite3

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/aodin/aspect"
)

func TestDatetime(t *testing.T) {
	expect := aspect.NewTester(t, &Sqlite3{})

	expect.Create("DATETIME", Datetime{})
	expect.Create(
		"DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP",
		Datetime{NotNull: true, Default: CurrentTimestamp},
	)

	d := time.Date(2014, 1, 1, 12, 0, 0, 0, time.UTC)
	value, err := Datetime{}.Validate(d)
	assert.Nil(t, err)
	assert.Equal(t, d, value)

	_, err = Datetime{}.Validate(123)
	assert.NotNil(t, err)
}
