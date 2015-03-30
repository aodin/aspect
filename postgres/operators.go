package postgres

import "github.com/aodin/aspect"

type Operator string

func (op Operator) String() string {
	return string(op)
}

func (op Operator) Compile(d aspect.Dialect, p *aspect.Parameters) (string, error) {
	return op.String(), nil
}

func (op Operator) With(name string) WithClause {
	return WithClause{Name: name, Operator: op}
}

const (
	Equal                     Operator = "="
	NotEqual                  Operator = "<>"
	LessThan                  Operator = "<"
	GreaterThan               Operator = ">"
	LessThanOrEqual           Operator = "<="
	GreaterThanOrEqual        Operator = ">="
	Contains                  Operator = "@>"
	ContainedBy               Operator = "<@"
	Overlap                   Operator = "&&"
	StrictlyLeftOf            Operator = "<<"
	StrictlyRightOf           Operator = ">>"
	DoesNotExtendToTheRightOf Operator = "&<"
	DoesNotExtendToTheLeftOf  Operator = "&>"
	IsAdjacentTo              Operator = "-|-"
	Union                     Operator = "+"
	Intersection              Operator = "*"
	Difference                Operator = "-"
)
