package postgres

import (
	"fmt"
	"github.com/aodin/aspect"
	"strings"
)

type Serial struct {
	PrimaryKey bool
	NotNull    bool
}

func (s Serial) Create(d aspect.Dialect) (string, error) {
	compiled := "SERIAL"
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

type Inet struct {
	NotNull    bool
	Unique     bool
	PrimaryKey bool
}

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
