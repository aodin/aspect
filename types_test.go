package aspect

import (
	"testing"
)

func TestBoolean(t *testing.T) {
	s := Boolean{}

	output, err := s.Create(&defaultDialect{})
	expected := "BOOL"
	if err != nil {
		t.Fatalf("Unexpected error during %s create: %s", expected, err)
	}
	if output != expected {
		t.Fatalf("Unexpected %s creation output: %s", expected, output)
	}

	s = Boolean{NotNull: true}
	output, err = s.Create(&defaultDialect{})
	expected = "BOOL NOT NULL"
	if err != nil {
		t.Fatalf("Unexpected error during %s create: %s", expected, err)
	}
	if output != expected {
		t.Fatalf("Unexpected %s creation output: %s", expected, output)
	}

	s = Boolean{NotNull: true, Default: "FALSE"}
	output, err = s.Create(&defaultDialect{})
	expected = "BOOL NOT NULL DEFAULT FALSE"
	if err != nil {
		t.Fatalf("Unexpected error during %s create: %s", expected, err)
	}
	if output != expected {
		t.Fatalf("Unexpected %s creation output: %s", expected, output)
	}
}
