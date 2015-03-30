package postgres

import (
	"testing"

	"github.com/aodin/aspect"
)

func TestOperators(t *testing.T) {
	expect := aspect.NewTester(t, &PostGres{})
	expect.SQL(`"room" WITH =`, Equal.With("room"))
}
