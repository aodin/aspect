package postgres

import (
	"testing"

	"github.com/aodin/aspect"
)

func TestUUID(t *testing.T) {
	expect := aspect.NewTester(t, &PostGres{})
	expect.Create("UUID", UUID{})
	expect.Create("UUID PRIMARY KEY", UUID{PrimaryKey: true})
	expect.Create(
		"UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4()",
		UUID{PrimaryKey: true, NotNull: true, Default: GenerateV4},
	)
}
