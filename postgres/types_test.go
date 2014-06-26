package postgres

import (
	"testing"
)

func TestSerial(t *testing.T) {
	s := Serial{}
	output, err := s.Create(&PostGres{})
	if err != nil {
		t.Fatalf("Unexpected error furing SERIAL create: %s", err)
	}
	expected := "SERIAL"
	if output != expected {
		t.Fatalf("Unexpected SERIAL creation output: %s", output)
	}

	s = Serial{PrimaryKey: true}
	output, err = s.Create(&PostGres{})
	if err != nil {
		t.Fatalf("Unexpected error furing SERIAL PRIMARY KEY create: %s", err)
	}
	expected = "SERIAL PRIMARY KEY"
	if output != expected {
		t.Fatalf("Unexpected SERIAL PRIMARY KEY creation output: %s", output)
	}
}
