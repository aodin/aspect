package postgres

import (
	"fmt"
	"strings"

	"github.com/aodin/aspect"
)

type JSON struct {
	PrimaryKey bool
	NotNull    bool
	Default    string // TODO clause that compiles?
}

var _ aspect.Type = JSON{}

func (s JSON) Create(d aspect.Dialect) (string, error) {
	compiled := "JSON"
	attrs := make([]string, 0)
	if s.PrimaryKey {
		attrs = append(attrs, "PRIMARY KEY")
	}
	if s.NotNull {
		attrs = append(attrs, "NOT NULL")
	}
	if len(attrs) > 0 {
		compiled += fmt.Sprintf(" %s", strings.Join(attrs, " "))
	}
	return compiled, nil
}

func (s JSON) IsPrimaryKey() bool {
	return s.PrimaryKey
}

func (s JSON) IsRequired() bool {
	return s.NotNull
}

func (s JSON) IsUnique() bool {
	return true
}

func (s JSON) Validate(i interface{}) (interface{}, error) {
	// TODO validation of JSON?
	return i, nil
}
