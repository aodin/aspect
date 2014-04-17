package aspect

import (
	"fmt"
	"strings"
)

// TODO Integration with clauses
type WrappedColumn struct {
	column ColumnElement
	f      string
	args   []string
	params []*Parameter
}

func (c *WrappedColumn) String() string {
	compiled, _ := c.Compile(&PostGres{}, Params())
	return compiled
}

func (c *WrappedColumn) Compile(d Dialect, params *Parameters) (string, error) {
	cc, err := c.column.Compile(d, params)
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

func (c *WrappedColumn) Name() string {
	return c.column.Name()
}

func (c *WrappedColumn) Table() *TableElem {
	return c.column.Table()
}

// Implement the Selectable interface for building SELECT statements
// This interface must be implemented for aggregate column or it will use
// the embedded ColumnStruct method.
func (c *WrappedColumn) Selectable() []ColumnElement {
	return []ColumnElement{c}
}

// Same with the Orderable interface
func (c *WrappedColumn) Orderable() *OrderedColumn {
	return &OrderedColumn{inner: c}
}

// TODO I have to re-implement all the ordering methods?
func (c *WrappedColumn) Asc() *OrderedColumn {
	return &OrderedColumn{inner: c}
}

func (c *WrappedColumn) Desc() *OrderedColumn {
	return &OrderedColumn{inner: c, desc: true}
}

func (c *WrappedColumn) NullsFirst() *OrderedColumn {
	return &OrderedColumn{inner: c, nullsFirst: true}
}

func (c *WrappedColumn) NullsLast() *OrderedColumn {
	return &OrderedColumn{inner: c, nullsLast: true}
}

func Avg(c *ColumnStruct) *WrappedColumn {
	return &WrappedColumn{column: c, f: "AVG"}
}

func Count(c *ColumnStruct) *WrappedColumn {
	return &WrappedColumn{column: c, f: "COUNT"}
}

func Date(c *ColumnStruct) *WrappedColumn {
	return &WrappedColumn{column: c, f: "DATE"}
}

func Max(c *ColumnStruct) *WrappedColumn {
	return &WrappedColumn{column: c, f: "MAX"}
}

func DatePart(c *ColumnStruct, part string) *WrappedColumn {
	// TODO add the part as a parameter?
	return &WrappedColumn{
		column: c,
		f:      "DATE_PART",
		args:   []string{fmt.Sprintf("'%s'", part)},
	}
}
