package aspect

import (
	"fmt"
	"strings"
	"time"
)

// Date represents DATE column types.
type Date struct {
	NotNull    bool
	PrimaryKey bool
	Unique     bool
}

var _ Type = Date{}

// Create returns the syntax need to create this column in CREATE statements.
func (s Date) Create(d Dialect) (string, error) {
	compiled := "DATE"
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

func (s Date) IsPrimaryKey() bool {
	return s.PrimaryKey
}

func (s Date) IsRequired() bool {
	return s.NotNull
}

func (s Date) IsUnique() bool {
	return s.PrimaryKey
}

func (s Date) Validate(i interface{}) (interface{}, error) {
	if _, ok := i.(time.Time); !ok {
		return i, fmt.Errorf("value is of non-time type %T", i)
	}
	return i, nil
}
