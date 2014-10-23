package aspect

import (
	"testing"
)

func TestPrimaryKey(t *testing.T) {
	// Test contains
	pk := PrimaryKeyArray{"id"}
	if !pk.Contains("id") {
		t.Errorf("pk should contain a column named 'id'")
	}
	if pk.Contains("dne") {
		t.Errorf("pk does not contain a column named 'dne'")
	}
}
