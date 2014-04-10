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

// Implement the sql.Selectable interface for building SELECT statements
// This interface must be implemented for aggregate column or it will use
// the embedded ColumnStruct method.
func (c *AggregateColumn) Selectable() []ColumnElement {
	return []ColumnElement{c}
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
