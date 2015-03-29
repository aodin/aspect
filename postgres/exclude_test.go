package postgres

import (
	"testing"

	"github.com/aodin/aspect"
)

func TestExclude(t *testing.T) {
	expect := aspect.NewTester(t, &PostGres{})
	expect.Create(
		`EXCLUDE ("room" WITH =)`,
		Exclude(Equal.With(times.C["room"])),
	)
	expect.Create(
		`EXCLUDE USING gist ("room" WITH =, "when" WITH &&)`,
		Exclude(
			Equal.With(times.C["room"]),
			Overlap.With(times.C["when"]),
		).Using(Gist),
	)
}
