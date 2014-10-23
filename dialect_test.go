package aspect

import (
	"testing"
)

func TestRegistry(t *testing.T) {
	// Add a dialect
	RegisterDialect("default", &defaultDialect{})

	// Get a dialect
	if _, err := GetDialect("default"); err != nil {
		t.Errorf("unexpected error getting a dialect that exists: %s", err)
	}

	// Get a dialect that doesn't exist
	if _, err := GetDialect("dne"); err == nil {
		t.Errorf("expected an error when getting a dialect that doesn't exist")
	}

	// Attempt to add a nil dialect
	func() {
		defer func() {
			if panicked := recover(); panicked == nil {
				t.Errorf("registry failed to panic when given a nil dialect")
			}
		}()
		RegisterDialect("nil", nil)
	}()

	// Attempt to add a duplicate dialect
	func() {
		defer func() {
			if panicked := recover(); panicked == nil {
				t.Errorf("registry failed to panic when given a duplicate dialect")
			}
		}()
		RegisterDialect("default", &defaultDialect{})
	}()
}
