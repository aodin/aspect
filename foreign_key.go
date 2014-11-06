package aspect

import (
	"fmt"
)

// fkType is an internal type representation. It implements the Creatable
// interface so it can be used in CREATE TABLE statements.
type fkType struct {
	name     string
	col      ColumnElem
	typ      Type
	onDelete *fkAction
	onUpdate *fkAction
}

var _ Createable = fkType{}

type fkAction string

// The following constants represent possible foreign key actions that can
// be used in ON DELETE and ON UPDATE clauses.
const (
	NoAction   fkAction = "NO ACTION"
	Restrict   fkAction = "RESTRICT"
	Cascade    fkAction = "CASCADE"
	SetNull    fkAction = "SET NULL"
	SetDefault fkAction = "SET DEFAULT"
)

// OnDelete adds an ON DELETE clause to the foreign key
func (fk fkType) OnDelete(b fkAction) fkType {
	fk.onDelete = &b
	return fk
}

// OnUpdate add an ON UPDATE clause to the foreign key
func (fk fkType) OnUpdate(b fkAction) fkType {
	fk.onUpdate = &b
	return fk
}

// Create returns the element's syntax for a CREATE TABLE statement.
func (fk fkType) Create(d Dialect) (string, error) {
	// Compile the type
	ct, err := fk.typ.Create(d)
	if err != nil {
		return "", err
	}
	compiled := fmt.Sprintf(
		`"%s" %s REFERENCES %s("%s")`,
		fk.name,
		ct,
		fk.col.table.Name,
		fk.col.name,
	)
	if fk.onDelete != nil {
		compiled += fmt.Sprintf(" ON DELETE %s", *fk.onDelete)
	}
	if fk.onUpdate != nil {
		compiled += fmt.Sprintf(" ON UPDATE %s", *fk.onUpdate)
	}
	return compiled, nil
}

// Modify implements the TableModifier interface. It creates a column and
// adds the same column to the create array.
func (fk fkType) Modify(t *TableElem) error {
	// No modifing nil table elements
	if t == nil {
		return fmt.Errorf("aspect: columns cannot modify a nil table")
	}

	// Column names cannot be blank
	// TODO Add more rules for column names
	if fk.name == "" {
		return fmt.Errorf("aspect: columns must have a name")
	}

	// Create the column for this table
	column := ColumnElem{
		inner: ColumnClause{name: fk.name, table: t},
		name:  fk.name,
		table: t,
		typ:   fk.typ,
	}

	// Add the column to the unique set of columns for this table
	if duplicate := t.C.Add(column); duplicate != nil {
		return duplicate
	}

	// Add the name to the table order
	t.order = append(t.order, column.name)

	// Add the fk to the create array
	t.creates = append(t.creates, fk)

	return nil
}

// ForeignKey creates a fkElem from the given name and column.
// The given column must already have an assigned table. The new ColumnElem
// will inherit its type from the given column's type.
// TODO Add the ability for self-referential foreign keys.
func ForeignKey(name string, fk ColumnElem, ts ...Type) fkType {
	if len(ts) > 1 {
		panic("aspect: foreign keys may only have one overriding type")
	}

	t := fk.typ

	// Allow the type to be overridden
	if len(ts) == 1 {
		t = ts[0]
	}
	return fkType{
		name: name,
		col:  fk,
		typ:  t,
	}
}
