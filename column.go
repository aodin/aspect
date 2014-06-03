package aspect

import (
	"fmt"
	"reflect"
)

/*

Column
------
Represents a single Column within a Table in the SQL spec.
Implements the `TableModifier` interface.

*/

// Maintains a unique set of columns
type ColumnSet map[string]ColumnElem

// Returns an error if a ColumnElem by this name already exists
// TODO Make this a private method?
func (set ColumnSet) Add(c ColumnElem) error {
	if _, exists := set[c.name]; exists {
		return fmt.Errorf("A column with the name %s already exists", c.name)
	}
	set[c.name] = c
	return nil
}

type ColumnElem struct {
	inner Clause
	name  string
	table *TableElem
	typ   dbType
}

func (c ColumnElem) String() string {
	compiled, _ := c.Compile(&PostGres{}, Params())
	return compiled
}

func (c ColumnElem) Inner() Clause {
	return c.inner
}

func (c ColumnElem) SetInner(clause Clause) ColumnElem {
	c.inner = clause
	return c
}

func (c ColumnElem) Compile(d Dialect, params *Parameters) (string, error) {
	if c.inner == nil {
		// Old behavior
		return fmt.Sprintf(`"%s"."%s"`, c.table.Name, c.name), nil
	} else {
		return c.inner.Compile(d, params)
	}
}

// TODO remove these
func (c ColumnElem) Name() string {
	return c.name
}

// TODO remove these
func (c ColumnElem) Table() *TableElem {
	return c.table
}

// Implement the sql.Selectable interface for building SELECT statements
func (c ColumnElem) Selectable() []ColumnElem {
	return []ColumnElem{c}
}

// Ordering
// --------

// Implement the sql.Orderable interface in order.go
func (c ColumnElem) Orderable() OrderedColumn {
	return OrderedColumn{inner: c}
}

// All other functions should return an OrderedColumn
func (c ColumnElem) Asc() OrderedColumn {
	return OrderedColumn{inner: c}
}

func (c ColumnElem) Desc() OrderedColumn {
	return OrderedColumn{inner: c, desc: true}
}

func (c ColumnElem) NullsFirst() OrderedColumn {
	return OrderedColumn{inner: c, nullsFirst: true}
}

func (c ColumnElem) NullsLast() OrderedColumn {
	return OrderedColumn{inner: c, nullsLast: true}
}

// Conditionals
// ------------

func (c ColumnElem) Equals(i interface{}) BinaryClause {
	return BinaryClause{
		Pre:  c,
		Post: &Parameter{i},
		Sep:  " = ",
	}
}

func (c ColumnElem) DoesNotEqual(i interface{}) BinaryClause {
	return BinaryClause{
		Pre:  c,
		Post: &Parameter{i},
		Sep:  " != ",
	}
}

func (c ColumnElem) LessThan(i interface{}) BinaryClause {
	return BinaryClause{
		Pre:  c,
		Post: &Parameter{i},
		Sep:  " < ",
	}
}

func (c ColumnElem) GreaterThan(i interface{}) BinaryClause {
	return BinaryClause{
		Pre:  c,
		Post: &Parameter{i},
		Sep:  " > ",
	}
}

func (c ColumnElem) LTE(i interface{}) BinaryClause {
	return BinaryClause{
		Pre:  c,
		Post: &Parameter{i},
		Sep:  " <= ",
	}
}

func (c ColumnElem) GTE(i interface{}) BinaryClause {
	return BinaryClause{
		Pre:  c,
		Post: &Parameter{i},
		Sep:  " >= ",
	}
}

// An interface is used because the args may be of any type: ints, strings...
// TODO an error if something other than a slice is added?
func (c ColumnElem) In(args interface{}) BinaryClause {
	// Create the inner array clause and parameters
	a := ArrayClause{Clauses: make([]Clause, 0), Sep: ", "}

	// Use reflect to get arguments from the interface only if it is a slice
	s := reflect.ValueOf(args)
	switch s.Kind() {
	case reflect.Slice:
		for i := 0; i < s.Len(); i++ {
			a.Clauses = append(a.Clauses, &Parameter{s.Index(i).Interface()})
		}
	}
	return BinaryClause{
		Pre:  c,
		Post: FuncClause{Inner: a},
		Sep:  " IN ",
	}
}

func (c ColumnElem) Between(a, b interface{}) ArrayClause {
	return AllOf(c.GTE(a), c.LTE(b))
}

func (c ColumnElem) NotBetween(a, b interface{}) ArrayClause {
	return AnyOf(c.LessThan(a), c.GreaterThan(b))
}

// Schema
// ------

// To implement the TableModifier interface the ColumnElem must
// have method Modify(). It does not need to modify its parent table.
func (c ColumnElem) Modify(t *TableElem) error {
	// No re-using columns across tables!
	if c.table != nil {
		return fmt.Errorf("Column %s already belongs to table %s", c.name, t.Name)
	}

	// Set the parent table of this column
	c.table = t

	// Update the inner clause with the completed ColumnClause
	c.inner = ColumnClause{table: t, name: c.name}

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
// TODO The constructor function does not need to return a ColumnElem,
// it can return a struct that modifies the table and adds a column.
func Column(name string, t dbType) ColumnElem {
	// Set the inner clause of the column to the incomplete ColumnClause.
	// This will be overwritten by the table modify function.
	return ColumnElem{inner: ColumnClause{name: name}, name: name, typ: t}
}
