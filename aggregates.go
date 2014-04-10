package aspect

import (
	"fmt"
)

type AggregateColumn struct {
	*ColumnStruct
	f string
}

func (c *AggregateColumn) String() string {
	return c.Compile()
}

func (c *AggregateColumn) Compile() string {
	return fmt.Sprintf(`%s(%s)`, c.f, c.ColumnStruct.Compile())
}

// Implement the Selectable interface for building SELECT statements
// This interface must be implemented for aggregate column or it will use
// the embedded ColumnStruct method.
func (c *AggregateColumn) Selectable() []ColumnElement {
	return []ColumnElement{c}
}

// Same with the Orderable interface
func (c *AggregateColumn) Orderable() *OrderedColumn {
	fmt.Println("Wrapping an aggregate column:", c.Compile())
	return &OrderedColumn{inner: c}
}

// TODO I have to re-implement all the ordering methods?
func (c *AggregateColumn) Asc() *OrderedColumn {
	return &OrderedColumn{inner: c}
}

func (c *AggregateColumn) Desc() *OrderedColumn {
	return &OrderedColumn{inner: c, desc: true}
}

func (c *AggregateColumn) NullsFirst() *OrderedColumn {
	return &OrderedColumn{inner: c, nullsFirst: true}
}

func (c *AggregateColumn) NullsLast() *OrderedColumn {
	return &OrderedColumn{inner: c, nullsLast: true}
}

func Count(column *ColumnStruct) *AggregateColumn {
	return &AggregateColumn{ColumnStruct: column, f: "COUNT"}
}

func Max(column *ColumnStruct) *AggregateColumn {
	return &AggregateColumn{ColumnStruct: column, f: "MAX"}
}

func Avg(column *ColumnStruct) *AggregateColumn {
	return &AggregateColumn{ColumnStruct: column, f: "AVG"}
}
