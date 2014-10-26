package aspect

import (
	"testing"
)

// TODO Generalize testing - see testify's assert

func TestText(t *testing.T) {
	s := Text{}
	output, err := s.Create(&defaultDialect{})
	expected := "TEXT"
	if err != nil {
		t.Fatalf("unexpected error during %s create: %s", expected, err)
	}
	if output != expected {
		t.Fatalf("expected %s, got %s", expected, output)
	}
}

func TestString(t *testing.T) {
	s := String{}
	output, err := s.Create(&defaultDialect{})
	expected := "VARCHAR"
	if err != nil {
		t.Fatalf("unexpected error during %s create: %s", expected, err)
	}
	if output != expected {
		t.Fatalf("expected %s, got %s", expected, output)
	}

	// Add attrs
	s = String{Length: 128, PrimaryKey: true, NotNull: true, Unique: true}
	output, err = s.Create(&defaultDialect{})
	expected = "VARCHAR(128) PRIMARY KEY NOT NULL UNIQUE"
	if err != nil {
		t.Fatalf("unexpected error during %s create: %s", expected, err)
	}
	if output != expected {
		t.Fatalf("expected %s, got %s", expected, output)
	}
}

func TestInteger(t *testing.T) {
	s := Integer{}
	output, err := s.Create(&defaultDialect{})
	expected := "INTEGER"
	if err != nil {
		t.Fatalf("unexpected error during %s create: %s", expected, err)
	}
	if output != expected {
		t.Fatalf("expected %s, got %s", expected, output)
	}

	// Add attrs
	s = Integer{PrimaryKey: true, NotNull: true, Unique: true}
	output, err = s.Create(&defaultDialect{})
	expected = "INTEGER PRIMARY KEY NOT NULL UNIQUE"
	if err != nil {
		t.Fatalf("unexpected error during %s create: %s", expected, err)
	}
	if output != expected {
		t.Fatalf("expected %s, got %s", expected, output)
	}
}

func TestTimestamp(t *testing.T) {
	s := Timestamp{Default: "now() at time zone 'utc'"}
	output, err := s.Create(&defaultDialect{})
	expected := "TIMESTAMP DEFAULT (now() at time zone 'utc')"
	if err != nil {
		t.Fatalf("unexpected error during %s create: %s", expected, err)
	}
	if output != expected {
		t.Fatalf("expected %s, got %s", expected, output)
	}

	s = Timestamp{
		WithTimezone: true,
		NotNull:      true,
		Default:      "now() at time zone 'utc'",
	}
	output, err = s.Create(&defaultDialect{})
	expected = "TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (now() at time zone 'utc')"
	if err != nil {
		t.Fatalf("unexpected error during %s create: %s", expected, err)
	}
	if output != expected {
		t.Fatalf("expected %s, got %s", expected, output)
	}
}

func TestDate(t *testing.T) {
	s := Date{PrimaryKey: true, NotNull: true, Unique: true}
	output, err := s.Create(&defaultDialect{})
	expected := "DATE PRIMARY KEY NOT NULL UNIQUE"
	if err != nil {
		t.Fatalf("unexpected error during %s create: %s", expected, err)
	}
	if output != expected {
		t.Fatalf("expected %s, got %s", expected, output)
	}
}

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

	s = Boolean{NotNull: true, Default: False}
	output, err = s.Create(&defaultDialect{})
	expected = "BOOL NOT NULL DEFAULT FALSE"
	if err != nil {
		t.Fatalf("Unexpected error during %s create: %s", expected, err)
	}
	if output != expected {
		t.Fatalf("Unexpected %s creation output: %s", expected, output)
	}
}

func TestDouble(t *testing.T) {
	s := Double{PrimaryKey: true, NotNull: true, Unique: true}
	output, err := s.Create(&defaultDialect{})
	expected := "DOUBLE PRECISION PRIMARY KEY NOT NULL UNIQUE"
	if err != nil {
		t.Fatalf("unexpected error during %s create: %s", expected, err)
	}
	if output != expected {
		t.Fatalf("expected %s, got %s", expected, output)
	}
}

func TestReal(t *testing.T) {
	s := Real{PrimaryKey: true, NotNull: true, Unique: true}
	output, err := s.Create(&defaultDialect{})
	expected := "REAL PRIMARY KEY NOT NULL UNIQUE"
	if err != nil {
		t.Fatalf("unexpected error during %s create: %s", expected, err)
	}
	if output != expected {
		t.Fatalf("expected %s, got %s", expected, output)
	}
}
