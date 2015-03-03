package aspect

import (
	"fmt"
	"time"
)

// Timestamp represents TIMESTAMP column types.
// TODO take a time.Location for timezone options
type Timestamp struct {
	NotNull         bool
	PrimaryKey      bool
	WithTimezone    bool
	WithoutTimezone bool
	Default         string // TODO Should this be a clause that compiles?
}

var _ Type = Timestamp{}

// Create returns the syntax need to create this column in CREATE statements.
func (s Timestamp) Create(d Dialect) (string, error) {
	compiled := "TIMESTAMP"
	if s.WithTimezone {
		compiled += " WITH TIME ZONE"
	} else if s.WithoutTimezone { // Only one timezone modifier is allowed
		compiled += " WITHOUT TIME ZONE"
	}
	if s.NotNull {
		compiled += " NOT NULL"
	}
	if s.Default != "" {
		compiled += fmt.Sprintf(" DEFAULT (%s)", s.Default)
	}
	return compiled, nil
}

func (s Timestamp) IsPrimaryKey() bool {
	return s.PrimaryKey
}

func (s Timestamp) IsRequired() bool {
	return s.NotNull && s.Default == ""
}

func (s Timestamp) IsUnique() bool {
	return s.PrimaryKey
}

func (s Timestamp) Validate(i interface{}) (interface{}, error) {
	if _, ok := i.(time.Time); !ok {
		return i, fmt.Errorf("value is of non-time type %T", i)
	}
	return i, nil
}
