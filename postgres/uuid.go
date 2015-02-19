package postgres

import (
	"fmt"
	"strings"

	"github.com/aodin/aspect"
)

var GenerateV4 string = `uuid_generate_v4()`

type UUID struct {
	PrimaryKey bool
	NotNull    bool
	Default    string // TODO clause that compiles?
}

var _ aspect.Type = UUID{}

func (s UUID) Create(d aspect.Dialect) (string, error) {
	compiled := "UUID"
	attrs := make([]string, 0)
	if s.PrimaryKey {
		attrs = append(attrs, "PRIMARY KEY")
	}
	if s.NotNull {
		attrs = append(attrs, "NOT NULL")
	}
	if s.Default != "" {
		attrs = append(attrs, fmt.Sprintf("DEFAULT %s", s.Default))
	}
	if len(attrs) > 0 {
		compiled += fmt.Sprintf(" %s", strings.Join(attrs, " "))
	}
	return compiled, nil
}

func (s UUID) IsPrimaryKey() bool {
	return s.PrimaryKey
}

func (s UUID) IsRequired() bool {
	return s.NotNull
}

func (s UUID) IsUnique() bool {
	return true
}

func (s UUID) Validate(i interface{}) (interface{}, error) {
	// TODO validation of UUID?
	return i, nil
}
