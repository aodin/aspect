package aspect

import (
	"fmt"
	"strings"
)

// All clauses must implement the Compiles interface
type Clause interface {
	Compiles
}

// Special clause type used for column selections
type ColumnClause struct {
	table *TableElem
	name  string
}

func (c ColumnClause) String() string {
	compiled, _ := c.Compile(&defaultDialect{}, Params())
	return compiled
}

func (c ColumnClause) Compile(d Dialect, params *Parameters) (string, error) {
	return fmt.Sprintf(`"%s"."%s"`, c.table.Name, c.name), nil
}

// TODO This is a dangerous clause that leads to parameters not being escaped
type StringClause struct {
	Name string
}

func (c StringClause) String() string {
	compiled, _ := c.Compile(&defaultDialect{}, Params())
	return compiled
}

func (c StringClause) Compile(d Dialect, params *Parameters) (string, error) {
	return fmt.Sprintf(`'%s'`, c.Name), nil
}

// TODO This is a dangerous clause that leads to parameters not being escaped
type IntClause struct {
	D int
}

func (c IntClause) String() string {
	compiled, _ := c.Compile(&defaultDialect{}, Params())
	return compiled
}

func (c IntClause) Compile(d Dialect, params *Parameters) (string, error) {
	return fmt.Sprintf(`%d`, c.D), nil
}

type FuncClause struct {
	Inner Clause
	F     string
}

func (c FuncClause) String() string {
	compiled, _ := c.Compile(&defaultDialect{}, Params())
	return compiled
}

func (c FuncClause) Compile(d Dialect, params *Parameters) (string, error) {
	cc, err := c.Inner.Compile(d, params)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s(%s)", c.F, cc), nil
}

type BinaryClause struct {
	Pre, Post Clause
	Sep       string
}

func (c BinaryClause) String() string {
	compiled, _ := c.Compile(&defaultDialect{}, Params())
	return compiled
}

func (c BinaryClause) Compile(d Dialect, params *Parameters) (string, error) {
	prec, err := c.Pre.Compile(d, params)
	if err != nil {
		return "", err
	}
	postc, err := c.Post.Compile(d, params)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s%s", prec, c.Sep, postc), nil
}

type ArrayClause struct {
	Clauses []Clause
	Sep     string
}

func (c ArrayClause) String() string {
	compiled, _ := c.Compile(&defaultDialect{}, Params())
	return compiled
}

func (c ArrayClause) Compile(d Dialect, params *Parameters) (string, error) {
	compiled := make([]string, len(c.Clauses))
	var err error
	for i, clause := range c.Clauses {
		compiled[i], err = clause.Compile(d, params)
		if err != nil {
			return "", err
		}
	}
	return strings.Join(compiled, c.Sep), nil
}

func AllOf(clauses ...Clause) ArrayClause {
	return ArrayClause{clauses, " AND "}
}

func AnyOf(clauses ...Clause) ArrayClause {
	return ArrayClause{clauses, " OR "}
}
