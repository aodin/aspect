package aspect

import (
	"fmt"
	"strings"
)

type dbType interface {
	Create(Dialect) (string, error)
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
	NotNull    bool
	PrimaryKey bool
}

func (s Timestamp) Create(d Dialect) (string, error) {
	compiled := "DATETIME"
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