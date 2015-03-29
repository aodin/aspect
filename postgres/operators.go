package postgres

import "github.com/aodin/aspect"

type Operator string

func (op Operator) String() string {
	return string(op)
}

func (op Operator) Compile(d aspect.Dialect, p *aspect.Parameters) (string, error) {
	return op.String(), nil
}

func (op Operator) With(column aspect.ColumnElem) ExcludeClause {
	return ExcludeClause{aspect.BinaryClause{
		Pre:  aspect.ColumnOnlyClause(column),
		Post: op,
		Sep:  " WITH ",
	}}
}

type ExcludeClause struct {
	aspect.BinaryClause
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
