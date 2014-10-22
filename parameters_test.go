package aspect

import (
	"testing"
)

func TestParameters(t *testing.T) {
	// Create a single parameter
	p := &Parameter{Value: "Hello"}
	output := p.String()
	expected := "$1"
	if output != expected {
		t.Errorf(
			"unexpected parameter String() output: %s != %s",
			expected,
			output,
		)
	}

	// Build multiple parameters
	ps := Params()
	ps.Add("Hello")
	if ps.Len() != 1 {
		t.Fatalf("unexpected length of params: 1 != %d", ps.Len())
	}

	// Get the args back out
	arg0 := ps.Args()[0]
	hello, ok := arg0.(string)
	if !ok {
		t.Fatalf("could not convert 'Hello' parameter to a string")
	}
	if hello != "Hello" {
		t.Errorf("unexpected string parameter: 'Hello' != %s", hello)
	}
}
