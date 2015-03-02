package aspect

import (
	"fmt"
	"log"
	"reflect"
)

// ColumnSet maintains a map of ColumnElem instances by column name
type ColumnSet map[string]ColumnElem

// Add adds a ColumnElem to the ColumnSet, returning an error if the name
// already exists.
func (set ColumnSet) Add(c ColumnElem) error {
	if _, exists := set[c.name]; exists {
		return fmt.Errorf(
			"aspect: column with the name %s already exists",
			c.name,
		)
	}
	set[c.name] = c
	return nil
}

// Column names must validate
func validateColumnName(name string) error {
	// TODO more rules
	if name == "" {
		return fmt.Errorf("aspect: column names cannot be blank")
	}
	return nil
}

// ColumnElem represents a table column. It implements the Compiles,
// Selectable, and Orderable interfaces for use in statements as well
// as the TableModifier and Creatable interfaces.
type ColumnElem struct {
	inner Clause
	name  string
	table *TableElem
	typ   Type
	alias string
}

var _ TableModifier = ColumnElem{}
var _ Creatable = ColumnElem{}

// clause generates a ColumnClause from this ColumnElem
func (c ColumnElem) clause() ColumnClause {
	return ColumnClause{table: c.table, name: c.name}
}

// As specifies an alias for this ColumnElem
func (c ColumnElem) As(alias string) ColumnElem {
	c.alias = alias
	return c
}

// String outputs a parameter-less SQL representation of the column using
// a neutral dialect. If an error occurred during compilation,
// then an empty string will be returned.
func (c ColumnElem) String() string {
	compiled, _ := c.Compile(&defaultDialect{}, Params())
	return compiled
}

// Inner returns the inner Clause of the ColumnElem. This operation is used
// for setting external casts and functions on the ColumnElem.
func (c ColumnElem) Inner() Clause {
	return c.inner
}

// SetInner sets the inner Clause of the ColumnElem. This operation is used
// for setting external casts and functions on the ColumnElem.
func (c ColumnElem) SetInner(clause Clause) ColumnElem {
	c.inner = clause
	return c
}

// Compile produces the dialect specific SQL and adds any parameters
// in the clause to the given Parameters instance
func (c ColumnElem) Compile(d Dialect, params *Parameters) (string, error) {
	var compiled string
	var err error
	if c.inner == nil {
		// Old behavior
		compiled, err = fmt.Sprintf(`"%s"."%s"`, c.table.Name, c.name), nil
	} else {
		compiled, err = c.inner.Compile(d, params)
	}
	if err != nil {
		return compiled, err
	}
	if c.alias != "" {
		compiled += fmt.Sprintf(` AS "%s"`, c.alias)
	}
	return compiled, nil
}

// Name returns the column's name
func (c ColumnElem) Name() string {
	return c.name
}

// Table returns the column's table
func (c ColumnElem) Table() *TableElem {
	return c.table
}

// Type returns the column's SQL type
func (c ColumnElem) Type() Type {
	return c.typ
}

// Selectable implements the sql.Selectable interface for building SELECT
// statements from ColumnElem instances.
func (c ColumnElem) Selectable() []ColumnElem {
	return []ColumnElem{c}
}

// Ordering
// --------

// Orerable implements the Orderable interface that allows the column itself
// to be used in an OrderBy clause.
func (c ColumnElem) Orderable() OrderedColumn {
	return OrderedColumn{inner: c}
}

// Asc returns an OrderedColumn. It is the same as passing the column itself
// to an OrderBy clause.
func (c ColumnElem) Asc() OrderedColumn {
	return OrderedColumn{inner: c}
}

// Desc returns an OrderedColumn that will sort in descending order.
func (c ColumnElem) Desc() OrderedColumn {
	return OrderedColumn{inner: c, desc: true}
}

// NullsFirst returns an OrderedColumn that will sort NULLs first.
func (c ColumnElem) NullsFirst() OrderedColumn {
	return OrderedColumn{inner: c, nullsFirst: true}
}

// NullsLast returns an OrderedColumn that will sort NULLs last.
func (c ColumnElem) NullsLast() OrderedColumn {
	return OrderedColumn{inner: c, nullsLast: true}
}

// Conditionals
// ------------

// Equals creates an equals clause that can be used in conditional clauses.
//  table.Select().Where(table.C["id"].Equals(3))
func (c ColumnElem) Equals(i interface{}) BinaryClause {
	// If the given parameter is a ColumnElem, wrap the post parameter
	// in a ColumnClause, not a Parameter
	switch t := i.(type) {
	case ColumnElem:
		return BinaryClause{
			Pre:  c,
			Post: t.clause(),
			Sep:  " = ",
		}
	default:
		return BinaryClause{
			Pre:  c,
			Post: &Parameter{i},
			Sep:  " = ",
		}
	}

}

// DoesNotEqual creates a does not equal clause that can be used in
// conditional clauses.
//  table.Select().Where(table.C["id"].DoesNotEqual(3))
func (c ColumnElem) DoesNotEqual(i interface{}) BinaryClause {
	return BinaryClause{
		Pre:  c,
		Post: &Parameter{i},
		Sep:  " != ",
	}
}

// LessThan creates a less than clause that can be used in conditional clauses.
//  table.Select().Where(table.C["id"].LessThan(3))
func (c ColumnElem) LessThan(i interface{}) BinaryClause {
	return BinaryClause{
		Pre:  c,
		Post: &Parameter{i},
		Sep:  " < ",
	}
}

// GreaterThan creates a greater than clause that can be used in conditional
// clauses.
//  table.Select().Where(table.C["id"].GreaterThan(3))
func (c ColumnElem) GreaterThan(i interface{}) BinaryClause {
	return BinaryClause{
		Pre:  c,
		Post: &Parameter{i},
		Sep:  " > ",
	}
}

// LTE creates a less than or equal to clause that can be used in conditional
// clauses.
//  table.Select().Where(table.C["id"].LTE(3))
func (c ColumnElem) LTE(i interface{}) BinaryClause {
	return BinaryClause{
		Pre:  c,
		Post: &Parameter{i},
		Sep:  " <= ",
	}
}

// GTE creates a greater than or equal to clause that can be used in
// conditional clauses.
//  table.Select().Where(table.C["id"].GTE(3))
func (c ColumnElem) GTE(i interface{}) BinaryClause {
	return BinaryClause{
		Pre:  c,
		Post: &Parameter{i},
		Sep:  " >= ",
	}
}

// Like creates a pattern matching clause that can be used in conditional
// clauses.
//  table.Select().Where(table.C["name"].Like(`_b%`))
func (c ColumnElem) Like(i string) BinaryClause {
	return BinaryClause{
		Pre:  c,
		Post: &Parameter{i},
		Sep:  " LIKE ",
	}
}

// NotLike creates a pattern matching clause that can be used in conditional
// clauses.
//  table.Select().Where(table.C["name"].NotLike(`_b%`))
func (c ColumnElem) NotLike(i string) BinaryClause {
	return BinaryClause{
		Pre:  c,
		Post: &Parameter{i},
		Sep:  " NOT LIKE ",
	}
}

// Like creates a case insensitive pattern matching clause that can be used in
// conditional clauses.
//  table.Select().Where(table.C["name"].ILike(`_b%`))
func (c ColumnElem) ILike(i string) BinaryClause {
	return BinaryClause{
		Pre:  c,
		Post: &Parameter{i},
		Sep:  " ILIKE ",
	}
}

// SimilarTo creates a SQL regular expression matching clause that can be used
// in conditional clauses.
//  table.Select().Where(table.C["name"].SimilarTo(`_b%`))
func (c ColumnElem) SimilarTo(i string) BinaryClause {
	return BinaryClause{
		Pre:  c,
		Post: &Parameter{i},
		Sep:  " SIMILAR TO ",
	}
}

// NotSimilarTo creates a SQL regular expression matching clause that can be
// used in conditional clauses.
//  table.Select().Where(table.C["name"].NotSimilarTo(`_b%`))
func (c ColumnElem) NotSimilarTo(i string) BinaryClause {
	return BinaryClause{
		Pre:  c,
		Post: &Parameter{i},
		Sep:  " NOT SIMILAR TO ",
	}
}

// IsNull creates a comparison clause that can be used for checking existence
// of NULLs in conditional clauses.
//  table.Select().Where(table.C["name"].IsNull())
func (c ColumnElem) IsNull() UnaryClause {
	return UnaryClause{
		Pre: c,
		Sep: " IS NULL",
	}
}

// IsNotNull creates a comparison clause that can be used for checking absence
// of NULLs in conditional clauses.
//  table.Select().Where(table.C["name"].IsNotNull())
func (c ColumnElem) IsNotNull() UnaryClause {
	return UnaryClause{
		Pre: c,
		Sep: " IS NOT NULL",
	}
}

// In creates a comparison clause with an IN operator that can be used in
// conditional clauses. An interface is used because the args may be of any
// type: ints, strings...
//  table.Select().Where(table.C["id"].In([]int64{1, 2, 3}))
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
	// TODO What if something other than a slice is given?
	// TODO This statement should be able to take clauses / subqueries
	return BinaryClause{
		Pre:  c,
		Post: FuncClause{Inner: a},
		Sep:  " IN ",
	}
}

func (c ColumnElem) Between(a, b interface{}) Clause {
	return AllOf(c.GTE(a), c.LTE(b))
}

func (c ColumnElem) NotBetween(a, b interface{}) Clause {
	return AnyOf(c.LessThan(a), c.GreaterThan(b))
}

// Schema
// ------

// Create implements the Creatable interface that outputs a column of a
// CREATE TABLE statement.
func (c ColumnElem) Create(d Dialect) (string, error) {
	ct, err := c.typ.Create(d)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`"%s" %s`, c.Name(), ct), nil
}

// Modify implements the TableModifier interface. It creates a column and
// adds the same column to the create array.
func (c ColumnElem) Modify(t *TableElem) error {
	// No modifing nil table elements
	if t == nil {
		return fmt.Errorf("aspect: columns cannot modify a nil table")
	}

	if err := validateColumnName(c.name); err != nil {
		return err
	}

	// No re-using columns across tables
	if c.table != nil {
		return fmt.Errorf(
			"aspect: column '%s' already belongs to table '%s'",
			c.name,
			t.Name,
		)
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
	t.order = append(t.order, c.name)

	// Add the column to the create array
	t.creates = append(t.creates, c)

	// If the type is a primary key, set the table primary key
	// However, if another pk is already set, panic
	if c.typ.IsPrimaryKey() {
		if t.pk != nil {
			log.Panicf(
				"aspect: cannot set column '%s' as PRIMARY KEY - there is already a primary key set: '%v' - try using a composite primary key with PrimaryKey()",
				c.name,
				t.pk,
			)
		}
		t.pk = PrimaryKeyArray{c.name}
	} else if c.typ.IsUnique() {
		// If the column type is unique, add it to the table's unique
		// constraints TODO Should primary keys beå added to uniques?
		t.uniques = append(t.uniques, UniqueConstraint{c.name})
	}
	return nil
}

// Column constructs a ColumnElem with the given name and Type.
func Column(name string, t Type) ColumnElem {
	// Set the inner clause of the column to the incomplete ColumnClause.
	// This will be overwritten by the table modify function.
	return ColumnElem{
		inner: ColumnClause{name: name},
		name:  name,
		typ:   t,
	}
}
