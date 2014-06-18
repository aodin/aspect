package aspect

import (
	"fmt"
	"strings"
)

type dbType interface {
	Create(Dialect) (string, error)
}

type Text struct{}

func (s Text) Create(d Dialect) (string, error) {
	compiled := "TEXT"
	return compiled, nil
}

type String struct {
	Length     int
	NotNull    bool
	Unique     bool
	PrimaryKey bool
}

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
	if len(attrs) > 0 {
		compiled += fmt.Sprintf(" %s", strings.Join(attrs, " "))
	}
	return compiled, nil
}

type Integer struct {
	NotNull    bool
	Unique     bool
	PrimaryKey bool
}

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

type Timestamp struct {
	NotNull      bool
	PrimaryKey   bool
	WithTimezone bool
}

func (s Timestamp) Create(d Dialect) (string, error) {
	compiled := "TIMESTAMP"
	if s.WithTimezone {
		compiled += " WITH TIME ZONE"
	}
	if s.NotNull {
		compiled += " NOT NULL"
	}
	return compiled, nil
}

type Date struct {
	NotNull    bool
	PrimaryKey bool
}

func (s Date) Create(d Dialect) (string, error) {
	compiled := "DATE"
	if s.NotNull {
		compiled += " NOT NULL"
	}
	return compiled, nil
}

type Boolean struct {
	NotNull bool
}

func (s Boolean) Create(d Dialect) (string, error) {
	compiled := "BOOL"
	if s.NotNull {
		compiled += " NOT NULL"
	}
	return compiled, nil
}

type Serial struct{}

func (s Serial) Create(d Dialect) (string, error) {
	return "SERIAL NOT NULL", nil
}

type Double struct {
	NotNull    bool
	Unique     bool
	PrimaryKey bool
}

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

type Real struct {
	NotNull    bool
	Unique     bool
	PrimaryKey bool
}

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

type Inet struct {
	NotNull    bool
	Unique     bool
	PrimaryKey bool
}

func (s Inet) Create(d Dialect) (string, error) {
	compiled := "INET"
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
