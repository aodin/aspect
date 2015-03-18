package postgres

import (
	"testing"

	"github.com/aodin/aspect"
)

func TestJSON(t *testing.T) {
	expect := aspect.NewTester(t, &PostGres{})
	expect.Create("JSON", JSON{})
	expect.Create("JSON PRIMARY KEY", JSON{PrimaryKey: true})
	expect.Create(
		"JSON PRIMARY KEY NOT NULL",
		JSON{PrimaryKey: true, NotNull: true},
	)
}
