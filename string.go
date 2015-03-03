package aspect

import (
	"fmt"
	"strings"
)

// Provide a blank string for default values
var (
	internalBlank string = ""
	Blank                = &internalBlank
)

// String represents VARCHAR column types.
type String struct {
	Length     int
	NotNull    bool
	Unique     bool
	PrimaryKey bool
	Default    *string
}

var _ Type = String{}

// Create returns the syntax need to create this column in CREATE statements.
func (s String) Create(d Dialect) (string, error) {
	compiled := "VARCHAR"
	attrs := make([]string, 0)
	if s.Length != 0 {
		compiled += fmt.Sprintf("(%d)", s.Length)
	}
	if s.PrimaryKey {
		attrs = append(attrs, "PRIMARY KEY")
	}
	if s.NotNull {
		attrs = append(attrs, "NOT NULL")
	}
	if s.Unique {
		attrs = append(attrs, "UNIQUE")
	}
	if s.Default != nil {
		attrs = append(attrs, fmt.Sprintf("DEFAULT '%s'", *s.Default))
	}
	if len(attrs) > 0 {
		compiled += fmt.Sprintf(" %s", strings.Join(attrs, " "))
	}
	return compiled, nil
}

func (s String) IsPrimaryKey() bool {
	return s.PrimaryKey
}

func (s String) IsRequired() bool {
	return s.NotNull
}

func (s String) IsUnique() bool {
	return s.PrimaryKey || s.Unique
}

func (s String) Validate(i interface{}) (interface{}, error) {
	value, ok := i.(string)
	if !ok {
		return i, fmt.Errorf("value is of non string type %T", i)
	}
	if s.Length != 0 && len(value) > s.Length {
		return i, fmt.Errorf(
			"value is %d characters long, max length is %d",
			len(value),
			s.Length,
		)
	}
	return i, nil
}
