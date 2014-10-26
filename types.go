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

type dbType interface {
	Create(Dialect) (string, error)
}

// Text represents TEXT column types.
type Text struct{}

// Create returns the syntax need to create this column in CREATE statements.
func (s Text) Create(d Dialect) (string, error) {
	compiled := "TEXT"
	return compiled, nil
}

// String represents VARCHAR column types.
type String struct {
	Length     int
	NotNull    bool
	Unique     bool
	PrimaryKey bool
}

// Create returns the syntax need to create this column in CREATE statements.
func (s String) Create(d Dialect) (string, error) {
	compiled := "VARCHAR"
	attrs := make([]string, 0)
	if s.Length != 0 {
		compiled += fmt.Sprintf("(%d)", s.Length)
	}
	// TODO Primary key implies unique
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

// Integer represents INTEGER column types.
type Integer struct {
	NotNull    bool
	Unique     bool
	PrimaryKey bool
}

// Create returns the syntax need to create this column in CREATE statements.
func (s Integer) Create(d Dialect) (string, error) {
	compiled := "INTEGER"
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

// Timestamp represents TIMESTAMP column types.
// TODO take a time.Location for timezone options
type Timestamp struct {
	NotNull      bool
	PrimaryKey   bool
	WithTimezone bool
	Default      string // TODO Should this be a clause that compiles?
}

// Create returns the syntax need to create this column in CREATE statements.
func (s Timestamp) Create(d Dialect) (string, error) {
	compiled := "TIMESTAMP"
	if s.WithTimezone {
		compiled += " WITH TIME ZONE"
	}
	if s.NotNull {
		compiled += " NOT NULL"
	}
	if s.Default != "" {
		compiled += fmt.Sprintf(" DEFAULT (%s)", s.Default)
	}
	return compiled, nil
}

// Date represents DATE column types.
type Date struct {
	NotNull    bool
	PrimaryKey bool
	Unique     bool
}

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

// Boolean represents BOOL column types. It contains a Default field that
// can be left nil, or set with the included variables True and False.
type Boolean struct {
	NotNull bool
	Default *bool
}

// Create returns the syntax need to create this column in CREATE statements.
func (s Boolean) Create(d Dialect) (string, error) {
	compiled := "BOOL"
	if s.NotNull {
		compiled += " NOT NULL"
	}
	if s.Default != nil {
		compiled += strings.ToUpper(fmt.Sprintf(" DEFAULT %t", *s.Default))
	}
	return compiled, nil
}

// Double represents DOUBLE PRECISION column types.
type Double struct {
	NotNull    bool
	Unique     bool
	PrimaryKey bool
}

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

// Real represents REAL column types.
type Real struct {
	NotNull    bool
	Unique     bool
	PrimaryKey bool
}

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
