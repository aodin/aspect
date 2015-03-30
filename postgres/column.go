package postgres

import "github.com/aodin/aspect"

// ColumnElem expands aspect's ColumnElem by adding PostGres specific operators
type ColumnElem struct {
	aspect.ColumnElem
}

func (c ColumnElem) operator(i interface{}, op Operator) aspect.BinaryClause {
	return aspect.BinaryClause{
		Pre:  c,
		Post: &aspect.Parameter{i},
		Sep:  " " + string(op) + " ",
	}
}

func (c ColumnElem) NotEqual(i interface{}) aspect.BinaryClause {
	return c.operator(i, NotEqual)
}

func (c ColumnElem) LessThan(i interface{}) aspect.BinaryClause {
	return c.operator(i, LessThan)
}

func (c ColumnElem) GreaterThan(i interface{}) aspect.BinaryClause {
	return c.operator(i, GreaterThan)
}

func (c ColumnElem) LessThanOrEqual(i interface{}) aspect.BinaryClause {
	return c.operator(i, LessThanOrEqual)
}

func (c ColumnElem) GreaterThanOrEqual(i interface{}) aspect.BinaryClause {
	return c.operator(i, GreaterThanOrEqual)
}

func (c ColumnElem) Contains(i interface{}) aspect.BinaryClause {
	return c.operator(i, Contains)
}

func (c ColumnElem) ContainedBy(i interface{}) aspect.BinaryClause {
	return c.operator(i, ContainedBy)
}

func (c ColumnElem) Overlap(i interface{}) aspect.BinaryClause {
	return c.operator(i, Overlap)
}

func (c ColumnElem) StrictlyLeftOf(i interface{}) aspect.BinaryClause {
	return c.operator(i, StrictlyLeftOf)
}

func (c ColumnElem) StrictlyRightOf(i interface{}) aspect.BinaryClause {
	return c.operator(i, StrictlyRightOf)
}

func (c ColumnElem) DoesNotExtendToTheRightOf(i interface{}) aspect.BinaryClause {
	return c.operator(i, DoesNotExtendToTheRightOf)
}

func (c ColumnElem) DoesNotExtendToTheLeftOf(i interface{}) aspect.BinaryClause {
	return c.operator(i, DoesNotExtendToTheLeftOf)
}

func (c ColumnElem) IsAdjacentTo(i interface{}) aspect.BinaryClause {
	return c.operator(i, IsAdjacentTo)
}

func (c ColumnElem) Union(i interface{}) aspect.BinaryClause {
	return c.operator(i, Union)
}

func (c ColumnElem) Intersection(i interface{}) aspect.BinaryClause {
	return c.operator(i, Intersection)
}

func (c ColumnElem) Difference(i interface{}) aspect.BinaryClause {
	return c.operator(i, Difference)
}

// Column wraps an aspect ColumnElem and adds postgres specific functionality
// C wraps an aspect ColumnElem and adds postgres specific functionality
func Column(column aspect.ColumnElem) ColumnElem {
	return C(column)
}

// C is a shorthand for Column
func C(column aspect.ColumnElem) ColumnElem {
	return ColumnElem{column}
}
