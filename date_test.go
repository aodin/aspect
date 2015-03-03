package aspect

import "testing"

func TestDate(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	expect.Create(
		"DATE PRIMARY KEY NOT NULL UNIQUE",
		Date{PrimaryKey: true, NotNull: true, Unique: true},
	)
}
