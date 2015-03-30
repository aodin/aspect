package postgres

import (
	"testing"

	"github.com/aodin/aspect"
)

func TestColumn(t *testing.T) {
	expect := aspect.NewTester(t, &PostGres{})
	expect.SQL(
		`"times"."when" @> $1`,
		C(times.C["when"]).Contains("[0, 1)"),
		"[0, 1)",
	)
}
