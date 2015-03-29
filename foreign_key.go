package aspect

import (
	"fmt"
	"log"
)

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

// ForeignKeyElem is an internal type representation. It implements the
// Creatable interface so it can be used in CREATE TABLE statements.
type ForeignKeyElem struct {
	name     string
	col      ColumnElem
	typ      Type
	table    *TableElem // the parent table of the key
	refTable *TableElem // the table the key references
	onDelete *fkAction
	onUpdate *fkAction
}

var _ Creatable = ForeignKeyElem{}

// Create returns the element's syntax for a CREATE TABLE statement.
func (fk ForeignKeyElem) Create(d Dialect) (string, error) {
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

func (fk ForeignKeyElem) ForeignName() string {
	return fk.col.name
}

// Modify implements the TableModifier interface. It creates a column and
// adds the same column to the create array.
func (fk ForeignKeyElem) Modify(t *TableElem) error {
	if t == nil {
		return fmt.Errorf("aspect: columns cannot modify a nil table")
	}

	if fk.table != nil {
		return fmt.Errorf(
			"aspect: foreign keys cannot be assigned to multiple tables",
		)
	}
	fk.table = t

	// Column names must validate
	if err := validateColumnName(fk.name); err != nil {
		return err
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

	// Add it to the list of foreign keys
	t.fks = append(t.fks, fk)

	return nil
}

func (fk ForeignKeyElem) Name() string {
	return fk.name
}

// OnDelete adds an ON DELETE clause to the foreign key
func (fk ForeignKeyElem) OnDelete(b fkAction) ForeignKeyElem {
	fk.onDelete = &b
	return fk
}

// OnUpdate add an ON UPDATE clause to the foreign key
func (fk ForeignKeyElem) OnUpdate(b fkAction) ForeignKeyElem {
	fk.onUpdate = &b
	return fk
}

// ReferencesTable returns the table that this foreign key references.
func (fk ForeignKeyElem) ReferencesTable() *TableElem {
	return fk.refTable
}

// Table returns the parent table of this foreign key.
func (fk ForeignKeyElem) Table() *TableElem {
	return fk.table
}

// Type returns the Type of this foreign key.
func (fk ForeignKeyElem) Type() Type {
	return fk.typ
}

// ForeignKey creates a ForeignKeyElem from the given name and column.
// The given column must already have an assigned table. The new ColumnElem
// will inherit its type from the given column's type, but a different
// type can be overridden by a single optional type.
func ForeignKey(name string, fk ColumnElem, ts ...Type) ForeignKeyElem {
	if fk.table == nil {
		log.Panic("aspect: foreign keys must reference a column with a table already assigned")
	}

	// Set the default type of the foreign key to the referencing column, but
	// allow the type to be overridden by a single optional type
	t := fk.typ
	if len(ts) > 1 {
		log.Panic("aspect: foreign keys may only have one overriding type")
	} else if len(ts) == 1 {
		t = ts[0]
	}
	return ForeignKeyElem{
		name:     name,
		col:      fk,
		typ:      t,
		refTable: fk.table,
	}
}

type SelfForeignKeyElem struct {
	ForeignKeyElem
	ref string
}

func (fk SelfForeignKeyElem) Modify(t *TableElem) error {
	if t == nil {
		return fmt.Errorf("aspect: columns cannot modify a nil table")
	}

	if fk.ForeignKeyElem.table != nil {
		return fmt.Errorf(
			"aspect: foreign keys cannot be assigned to multiple tables",
		)
	}

	// The ref column must also exist on the column
	if _, exists := t.C[fk.ref]; !exists {
		return fmt.Errorf(
			"aspect: no column with the name %s exists in the table %s",
			fk.ref,
			t.Name,
		)
	}

	fk.ForeignKeyElem.col = t.C[fk.ref]
	fk.ForeignKeyElem.table = t
	fk.ForeignKeyElem.refTable = t

	// If the type of fk is nil, use the column's type
	if fk.ForeignKeyElem.typ == nil {
		fk.ForeignKeyElem.typ = fk.ForeignKeyElem.col.typ
	}

	// Column names must validate
	if err := validateColumnName(fk.ForeignKeyElem.name); err != nil {
		return err
	}

	// Create the column for this table
	column := ColumnElem{
		inner: ColumnClause{name: fk.ForeignKeyElem.name, table: t},
		name:  fk.ForeignKeyElem.name,
		table: t,
		typ:   fk.ForeignKeyElem.typ,
	}

	// Add the column to the unique set of columns for this table
	if duplicate := t.C.Add(column); duplicate != nil {
		return duplicate
	}

	// Add the name to the table order
	t.order = append(t.order, column.name)

	// Add the fk to the create array
	t.creates = append(t.creates, fk)

	// Add it to the list of foreign keys
	t.fks = append(t.fks, fk.ForeignKeyElem)

	return nil
}

func SelfForeignKey(name, ref string, ts ...Type) SelfForeignKeyElem {
	// Set the default type of the foreign key to the referencing column, but
	// allow the type to be overridden by a single optional type
	var t Type
	if len(ts) > 1 {
		log.Panic("aspect: foreign keys may only have one overriding type")
	} else if len(ts) == 1 {
		t = ts[0]
	}
	return SelfForeignKeyElem{
		ForeignKeyElem: ForeignKeyElem{
			name: name,
			typ:  t,
		},
		ref: ref,
	}
}
