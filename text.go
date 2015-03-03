package aspect

import "fmt"

// Text represents TEXT column types.
type Text struct {
	NotNull bool
}

var _ Type = Text{}

// Create returns the syntax need to create this column in CREATE statements.
func (s Text) Create(d Dialect) (string, error) {
	compiled := "TEXT"
	if s.NotNull {
		compiled += " NOT NULL"
	}
	return compiled, nil
}

func (s Text) IsPrimaryKey() bool {
	return false
}

func (s Text) IsRequired() bool {
	return s.NotNull
}

func (s Text) IsUnique() bool {
	return false
}

func (s Text) Validate(i interface{}) (interface{}, error) {
	if _, ok := i.(string); !ok {
		return i, fmt.Errorf("value is of non string type %T", i)
	}
	return i, nil
}
