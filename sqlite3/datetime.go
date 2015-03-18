package sqlite3

import (
	"fmt"
	"time"

	"github.com/aodin/aspect"
)

// CurrentTimestamp is a default option for sqlite3 column types, but it
// only has second resolution
const CurrentTimestamp = `CURRENT_TIMESTAMP`

type Datetime struct {
	PrimaryKey bool
	NotNull    bool
	Default    string // TODO Should this be a clause that compiles?
}

var _ aspect.Type = Datetime{}

// Create returns the syntax need to create this column in CREATE statements.
func (s Datetime) Create(d aspect.Dialect) (string, error) {
	compiled := "DATETIME"
	if s.NotNull {
		compiled += " NOT NULL"
	}
	if s.Default != "" {
		compiled += fmt.Sprintf(" DEFAULT %s", s.Default)
	}
	return compiled, nil
}

func (s Datetime) IsPrimaryKey() bool {
	return s.PrimaryKey
}

func (s Datetime) IsRequired() bool {
	return s.NotNull
}

func (s Datetime) IsUnique() bool {
	return true
}

func (s Datetime) Validate(i interface{}) (interface{}, error) {
	if _, ok := i.(time.Time); !ok {
		return i, fmt.Errorf("value is of non-time type %T", i)
	}
	return i, nil
}
