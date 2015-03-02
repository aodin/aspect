package aspect

import (
	"fmt"
	"strings"
)

// Clause is a common interface that most components implement.
type Clause interface {
	Compiles
}

// ColumnClause is a clause used for column selections.
type ColumnClause struct {
	table *TableElem
	name  string
}

// String returns the ColumnClause's SQL using the default dialect.
func (c ColumnClause) String() string {
	compiled, _ := c.Compile(&defaultDialect{}, Params())
	return compiled
}

// Compile creates the SQL to represent a table column using the given
// dialect, optionally without a table prefix.
func (c ColumnClause) Compile(d Dialect, params *Parameters) (string, error) {
	if c.table == nil {
		return fmt.Sprintf(`"%s"`, c.name), nil
	} else {
		return fmt.Sprintf(`"%s"."%s"`, c.table.Name, c.name), nil
	}
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

type UnaryClause struct {
	Pre Clause
	Sep string
}

func (c UnaryClause) String() string {
	compiled, _ := c.Compile(&defaultDialect{}, Params())
	return compiled
}

func (c UnaryClause) Compile(d Dialect, params *Parameters) (string, error) {
	var pre string
	var err error
	if c.Pre != nil {
		pre, err = c.Pre.Compile(d, params)
		if err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("%s%s", pre, c.Sep), nil
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
	var pre, post string
	var err error
	if c.Pre != nil {
		pre, err = c.Pre.Compile(d, params)
		if err != nil {
			return "", err
		}
	}
	if c.Post != nil {
		post, err = c.Post.Compile(d, params)
		if err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("%s%s%s", pre, c.Sep, post), nil
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

// AllOf joins the given clauses with 'AND' and wraps them in parentheses
func AllOf(clauses ...Clause) Clause {
	return FuncClause{Inner: ArrayClause{clauses, " AND "}}
}

// AnyOf joins the given clauses with 'OR' and wraps them in parentheses
func AnyOf(clauses ...Clause) Clause {
	return FuncClause{Inner: ArrayClause{clauses, " OR "}}
}
