package postgres

import (
	"fmt"
	"strings"

	"github.com/aodin/aspect"
)

type Inet struct {
	NotNull    bool
	Unique     bool
	PrimaryKey bool
}

var _ aspect.Type = Inet{}

// Create returns the syntax need to create this column in CREATE statements.
func (s Inet) Create(d aspect.Dialect) (string, error) {
	compiled := "INET"
	attrs := make([]string, 0)
	if s.PrimaryKey {
		attrs = append(attrs, "PRIMARY KEY")
	}
	if s.NotNull {
		attrs = append(attrs, "NOT NULL")
	}
	if s.Unique {
		attrs = append(attrs, "UNIQUE")
	}
	if len(attrs) > 0 {
		compiled += fmt.Sprintf(" %s", strings.Join(attrs, " "))
	}
	return compiled, nil
}

func (s Inet) IsPrimaryKey() bool {
	return s.PrimaryKey
}

func (s Inet) IsRequired() bool {
	return s.NotNull
}

func (s Inet) IsUnique() bool {
	return s.PrimaryKey || s.Unique
}

func (s Inet) Validate(i interface{}) (interface{}, error) {
	// TODO validation of Inet?
	return i, nil
}
