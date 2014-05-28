package aspect

import (
	"testing"
	"time"
)

func TestColumn(t *testing.T) {
	// Load the Denver timezone
	denver, err := time.LoadLocation("America/Denver")
	if err != nil {
		t.Fatal(err)
	}
	expectedPostGres(
		t,
		DateOf(views.C["timestamp"].InLocation(denver)),
		`DATE("views"."timestamp"::TIMESTAMP WITH TIME ZONE AT TIME ZONE $1)`,
		1,
	)
}