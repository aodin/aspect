package aspect

import (
	"fmt"
	"strings"
)

// Provide a nullable boolean for Boolean type default values
var (
	internalTrue  bool = true
	internalFalse bool = false
	True               = &internalTrue
	False              = &internalFalse
)

// Boolean represents BOOL column types. It contains a Default field that
// can be left nil, or set with the included variables True and False.
type Boolean struct {
	NotNull bool
	Default *bool
}

var _ Type = Boolean{}

// Create returns the syntax need to create this column in CREATE statements.
func (s Boolean) Create(d Dialect) (string, error) {
	compiled := "BOOLEAN"
	if s.NotNull {
		compiled += " NOT NULL"
	}
	if s.Default != nil {
		compiled += strings.ToUpper(fmt.Sprintf(" DEFAULT %t", *s.Default))
	}
	return compiled, nil
}

func (s Boolean) IsPrimaryKey() bool {
	return false
}

func (s Boolean) IsRequired() bool {
	return s.NotNull && s.Default == nil
}

func (s Boolean) IsUnique() bool {
	return false
}

func (s Boolean) Validate(i interface{}) (interface{}, error) {
	// TODO parse boolean strings
	if _, ok := i.(bool); !ok {
		return i, fmt.Errorf("value is of non-bool type %T", i)
	}
	return i, nil
}
