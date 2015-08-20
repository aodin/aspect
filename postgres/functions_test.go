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

func Test_StringAgg(t *testing.T) {
	expect := aspect.NewTester(t, &PostGres{})
	expect.SQL(
		`string_agg("times"."when", ', ')`,
		StringAgg(times.C["when"], ", "),
	)
}
