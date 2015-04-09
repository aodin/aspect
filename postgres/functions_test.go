package postgres

import (
	"testing"

	"github.com/aodin/aspect"
)

func TestFunctions(t *testing.T) {
	expect := aspect.NewTester(t, &PostGres{})
	expect.SQL(
		`ISEMPTY("times"."when" * $1)`,
		IsEmpty(C(times.C["when"]).Intersection("range")),
		"range",
	)
}
