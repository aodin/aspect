package aspect

import (
	"fmt"
)

/*

Column
------
Represents a single Column within a Table in the SQL spec.
Implements the `TableModifier` interface.

*/

type ColumnElement interface {
	Clause
	Name() string
	Table() *TableElem
}

// Maintains a unique set of columns
type ColumnSet map[string]ColumnStruct

// Returns an error if a ColumnStruct by this name already exists
// TODO Make this a private method?
func (set ColumnSet) Add(c ColumnStruct) error {
	if _, exists := set[c.name]; exists {
		return fmt.Errorf("A column with the name %s already exists", c.name)
	}
	set[c.name] = c
	return nil
}

type ColumnStruct struct {
	name  string
	table *TableElem
	typ   dbType
}

func (c ColumnStruct) String() string {
	compiled, _ := c.Compile(&PostGres{}, Params())
	return compiled
}

func (c ColumnStruct) Compile(d Dialect, params *Parameters) (string, error) {
	return fmt.Sprintf(`"%s"."%s"`, c.table.Name, c.name), nil
}

func (c ColumnStruct) Name() string {
	return c.name
}

func (c ColumnStruct) Table() *TableElem {
	return c.table
}

// Implement the sql.Selectable interface for building SELECT statements
func (c ColumnStruct) Selectable() []ColumnElement {
	return []ColumnElement{c}
}

// Ordering
// --------

// Implement the sql.Orderable interface in order.go
func (c ColumnStruct) Orderable() OrderedColumn {
	return OrderedColumn{inner: c}
}

// All other functions should return an OrderedColumn
func (c ColumnStruct) Asc() OrderedColumn {
	return OrderedColumn{inner: c}
}

func (c ColumnStruct) Desc() OrderedColumn {
	return OrderedColumn{inner: c, desc: true}
}

func (c ColumnStruct) NullsFirst() OrderedColumn {
	return OrderedColumn{inner: c, nullsFirst: true}
}

func (c ColumnStruct) NullsLast() OrderedColumn {
	return OrderedColumn{inner: c, nullsLast: true}
}

// Conditionals
// ------------

func (c ColumnStruct) Equals(i interface{}) *BinaryClause {
	return &BinaryClause{
		pre:  c,
		post: &Parameter{i},
		sep:  "=",
	}
}

func (c ColumnStruct) LessThan(i interface{}) *BinaryClause {
	return &BinaryClause{
		pre:  c,
		post: &Parameter{i},
		sep:  "<",
	}
}

func (c ColumnStruct) GreaterThan(i interface{}) *BinaryClause {
	return &BinaryClause{
		pre:  c,
		post: &Parameter{i},
		sep:  ">",
	}
}

func (c ColumnStruct) LTE(i interface{}) *BinaryClause {
	return &BinaryClause{
		pre:  c,
		post: &Parameter{i},
		sep:  "<=",
	}
}

func (c ColumnStruct) GTE(i interface{}) *BinaryClause {
	return &BinaryClause{
		pre:  c,
		post: &Parameter{i},
		sep:  ">=",
	}
}

func (c ColumnStruct) Between(a, b interface{}) *ArrayClause {
	return AllOf(c.GTE(a), c.LTE(b))
}

func (c ColumnStruct) NotBetween(a, b interface{}) *ArrayClause {
	return AnyOf(c.LessThan(a), c.GreaterThan(b))
}

// Schema
// ------

// To implement the TableModifier interface the ColumnStruct must
// have method Modify(). It does not need to modify its parent table.
func (c ColumnStruct) Modify(t *TableElem) error {
	// No re-using columns across tables!
	if c.table != nil {
		return fmt.Errorf("Column %s already belongs to table %s", c.name, t.Name)
	}

	// Set the parent table of this column
	c.table = t

	// Add the column to the unique set of columns for this table
	if duplicate := t.C.Add(c); duplicate != nil {
		return duplicate
	}

	// Add the name to the table order
	// TODO Something other than an append operation
	t.order = append(t.order, c.name)
	return nil
}

// Constructor function
// TODO The constructor function does not need to return a ColumnStruct,
// it can return a struct that modifies the table and adds a column.
func Column(name string, t dbType) ColumnStruct {
	return ColumnStruct{name: name, typ: t}
}
