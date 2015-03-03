package aspect

import "testing"

func TestParameters(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	// Create a single parameter
	expect.SQL("$1", &Parameter{Value: "hello"}, "hello")
}
