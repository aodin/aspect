package aspect

import (
	"fmt"
	"strconv"
	"strings"
)

// Double represents DOUBLE PRECISION column types.
type Double struct {
	NotNull    bool
	Unique     bool
	PrimaryKey bool
}

var _ Type = Double{}

// Create returns the syntax need to create this column in CREATE statements.
func (s Double) Create(d Dialect) (string, error) {
	compiled := "DOUBLE PRECISION"
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

func (s Double) IsPrimaryKey() bool {
	return s.PrimaryKey
}

func (s Double) IsRequired() bool {
	return s.NotNull
}

func (s Double) IsUnique() bool {
	return s.PrimaryKey || s.Unique
}

func (s Double) Validate(i interface{}) (interface{}, error) {
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
