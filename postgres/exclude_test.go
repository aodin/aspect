package postgres

import (
	"testing"

	"github.com/aodin/aspect"
)

func TestExclude(t *testing.T) {
	expect := aspect.NewTester(t, &PostGres{})
	expect.Create(
		`EXCLUDE ("room" WITH =)`,
		Exclude(Equal.With("room")),
	)
	expect.Create(
		`EXCLUDE USING gist ("room" WITH =, "when" WITH &&)`,
		Exclude(Equal.With("room"), Overlap.With("when")).Using(Gist),
	)
}
