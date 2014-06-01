package aspect

import (
	"fmt"
	"strings"
)

// All clauses must implement the Compiles interface
type Clause interface {
	Compiles
}

type StringClause struct {
	name string
}

func (c StringClause) String() string {
	compiled, _ := c.Compile(&PostGres{}, Params())
	return compiled
}

func (c StringClause) Compile(d Dialect, params *Parameters) (string, error) {
	return fmt.Sprintf(`'%s'`, c.name), nil
}

// Special clause type used for column selections
type ColumnClause struct {
	table *TableElem
	name  string
}

func (c ColumnClause) String() string {
	compiled, _ := c.Compile(&PostGres{}, Params())
	return compiled
}

func (c ColumnClause) Compile(d Dialect, params *Parameters) (string, error) {
	return fmt.Sprintf(`"%s"."%s"`, c.table.Name, c.name), nil
}

type FuncClause struct {
	clause Clause
	f      string
}

func (c FuncClause) String() string {
	compiled, _ := c.Compile(&PostGres{}, Params())
	return compiled
}

func (c FuncClause) Compile(d Dialect, params *Parameters) (string, error) {
	cc, err := c.clause.Compile(d, params)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s(%s)", c.f, cc), nil
}

func MakeFuncClause(f string, clause Clause) FuncClause {
	return FuncClause{f: f, clause: clause}
}

type BinaryClause struct {
	pre, post Clause
	sep       string
}

func (c BinaryClause) String() string {
	compiled, _ := c.Compile(&PostGres{}, Params())
	return compiled
}

func (c BinaryClause) Compile(d Dialect, params *Parameters) (string, error) {
	prec, err := c.pre.Compile(d, params)
	if err != nil {
		return "", err
	}
	postc, err := c.post.Compile(d, params)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s%s", prec, c.sep, postc), nil
}

type ArrayClause struct {
	clauses []Clause
	sep     string
}

func (c ArrayClause) String() string {
	compiled, _ := c.Compile(&PostGres{}, Params())
	return compiled
}

func (c ArrayClause) Compile(d Dialect, params *Parameters) (string, error) {
	compiled := make([]string, len(c.clauses))
	var err error
	for i, clause := range c.clauses {
		compiled[i], err = clause.Compile(d, params)
		if err != nil {
			return "", err
		}
	}
	return strings.Join(compiled, c.sep), nil
}

func MakeArrayClause(sep string, clauses []Clause) ArrayClause {
	return ArrayClause{sep: sep, clauses: clauses}
}

func AllOf(clauses ...Clause) ArrayClause {
	return ArrayClause{clauses, " AND "}
}

func AnyOf(clauses ...Clause) ArrayClause {
	return ArrayClause{clauses, " OR "}
}
