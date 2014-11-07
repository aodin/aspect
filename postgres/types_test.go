package postgres

import (
	"testing"

	"github.com/aodin/aspect"
)

func TestSerial(t *testing.T) {
	expect := aspect.NewTester(t, &PostGres{})

	expect.Create("SERIAL", Serial{})
	expect.Create("SERIAL PRIMARY KEY", Serial{PrimaryKey: true})
	expect.Create(
		"SERIAL PRIMARY KEY NOT NULL",
		Serial{PrimaryKey: true, NotNull: true},
	)
}

func TestInet(t *testing.T) {
	expect := aspect.NewTester(t, &PostGres{})

	expect.Create("INET", Inet{})
}
