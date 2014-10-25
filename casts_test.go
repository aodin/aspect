package aspect

import (
	"testing"
	"time"
)

func TestColumn(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	// Load the Denver timezone
	denver, err := time.LoadLocation("America/Denver")
	if err != nil {
		t.Fatal(err)
	}

	expect.SQL(
		`DATE("views"."timestamp"::TIMESTAMP WITH TIME ZONE AT TIME ZONE $1)`,
		DateOf(views.C["timestamp"].InLocation(denver)),
		denver.String(),
	)
}
