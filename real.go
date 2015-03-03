package aspect

import (
	"fmt"
	"strconv"
	"strings"
)

// Real represents REAL column types.
type Real struct {
	NotNull    bool
	Unique     bool
	PrimaryKey bool
}

var _ Type = Real{}

// Create returns the syntax need to create this column in CREATE statements.
func (s Real) Create(d Dialect) (string, error) {
	compiled := "REAL"
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

func (s Real) IsPrimaryKey() bool {
	return s.PrimaryKey
}

func (s Real) IsRequired() bool {
	return s.NotNull
}

func (s Real) IsUnique() bool {
	return s.PrimaryKey || s.Unique
}

func (s Real) Validate(i interface{}) (interface{}, error) {
	switch t := i.(type) {
	case string:
		v, err := strconv.ParseFloat(t, 64)
		if err != nil {
			return i, err
		}
		i = v
	case float32:
	case float64:
	case int:
		i = float64(t)
	case int64:
		i = float64(t)
	default:
		return i, fmt.Errorf("value is non-numeric type %T", t)
	}
	return i, nil
}
