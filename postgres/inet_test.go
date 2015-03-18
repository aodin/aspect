package postgres

import (
	"testing"

	"github.com/aodin/aspect"
)

func TestInet(t *testing.T) {
	expect := aspect.NewTester(t, &PostGres{})

	expect.Create("INET", Inet{})
}
