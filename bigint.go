package aspect

import (
	"fmt"
	"strconv"
	"strings"
)

// BigInt represents BIGINT column types.
type BigInt struct {
	NotNull    bool
	Unique     bool
	PrimaryKey bool
}

var _ Type = BigInt{}

// Create returns the syntax need to create this column in CREATE statements.
func (s BigInt) Create(d Dialect) (string, error) {
	compiled := "BIGINT"
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

func (s BigInt) IsPrimaryKey() bool {
	return s.PrimaryKey
}

func (s BigInt) IsRequired() bool {
	return s.NotNull
}

func (s BigInt) IsUnique() bool {
	return s.PrimaryKey || s.Unique
}

func (s BigInt) Validate(i interface{}) (interface{}, error) {
	switch t := i.(type) {
	case string:
		v, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return i, err
		}
		i = v
	case float64:
		v := int64(t)
		if t != float64(v) {
			return i, fmt.Errorf(
				"value is numeric, but not a whole number: %f",
				t,
			)
		}
		i = v
	case int:
		i = int64(i.(int))
	case int64:
	default:
		return i, fmt.Errorf("value is non-numeric type %T", t)
	}
	return i, nil
}
