package aspect

import (
	"fmt"
	"strings"
)

// All clauses must implement the Compiler interface
type Clause interface {
	Compiler
}

// TODO Or just type Parameter []interface{}
type Parameter struct {
	Value interface{}
}

func (p *Parameter) String() string {
	compiled, _ := p.Compile(&PostGres{}, Params())
	return compiled
}

// Parameter compilation is dialect dependent. Some dialects, such as
// PostGres, also require the parameter index
func (p *Parameter) Compile(d Dialect, params *Parameters) (string, error) {
	i := params.Add(p.Value)
	return d.Parameterize(i), nil
}

type UnaryClause struct {
	clause Clause
	oper   string
	pre    bool
}

func (c *UnaryClause) String() string {
	compiled, _ := c.Compile(&PostGres{}, Params())
	return compiled
}

func (c *UnaryClause) Compile(d Dialect, params *Parameters) (string, error) {
	cc, err := c.clause.Compile(d, params)
	if err != nil {
		return "", err
	}
	if c.pre {
		return fmt.Sprintf("%s(%s)", c.oper, cc), nil
	} else {
		return fmt.Sprintf("(%s)%s", cc, c.oper), nil
	}
}

type BinaryClause struct {
	pre, post Clause
	sep       string
}

func (c *BinaryClause) String() string {
	compiled, _ := c.Compile(&PostGres{}, Params())
	return compiled
}

func (c *BinaryClause) Compile(d Dialect, params *Parameters) (string, error) {
	prec, err := c.pre.Compile(d, params)
	if err != nil {
		return "", err
	}
	postc, err := c.post.Compile(d, params)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s %s %s", prec, c.sep, postc), nil
}

type ArrayClause struct {
	clauses []Clause
	sep     string
}

func (c *ArrayClause) String() string {
	compiled, _ := c.Compile(&PostGres{}, Params())
	return compiled
}

func (c *ArrayClause) Compile(d Dialect, params *Parameters) (string, error) {
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

func AllOf(clauses ...Clause) *ArrayClause {
	return &ArrayClause{clauses, " AND "}
}

func AnyOf(clauses ...Clause) *ArrayClause {
	return &ArrayClause{clauses, " OR "}
}