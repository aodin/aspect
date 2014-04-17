package aspect

import (
	"fmt"
	"strings"
)

// TODO Another implementation for functions?
// TODO args are pre or post?
type AggregateColumn struct {
	*ColumnStruct
	f    string
	args []string
}

func (c *AggregateColumn) String() string {
	compiled, _ := c.Compile(&PostGres{}, Params())
	return compiled
}

func (c *AggregateColumn) Compile(d Dialect, params *Parameters) (string, error) {
	cc, err := c.ColumnStruct.Compile(d, params)
	if err != nil {
		return "", err
	}
	if c.args != nil && len(c.args) > 0 {
		// TODO parameterization of args
		return fmt.Sprintf(
			`%s(%s, %s)`,
			c.f,
			strings.Join(c.args, ", "),
			cc,
		), nil
	}
	return fmt.Sprintf(`%s(%s)`, c.f, cc), nil
}

// Implement the Selectable interface for building SELECT statements
// This interface must be implemented for aggregate column or it will use
// the embedded ColumnStruct method.
func (c *AggregateColumn) Selectable() []ColumnElement {
	return []ColumnElement{c}
}

// Same with the Orderable interface
func (c *AggregateColumn) Orderable() *OrderedColumn {
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

func DatePart(column *ColumnStruct, part string) *AggregateColumn {
	// TODO add the part as a parameter?
	return &AggregateColumn{
		ColumnStruct: column,
		f:            "DATE_PART",
		args:         []string{fmt.Sprintf("'%s'", part)},
	}
}
