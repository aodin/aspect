package aspect

import (
	"fmt"
	"strings"
)

// fkType is an internal type representation
type fkType struct {
	t dbType
	c ColumnElem
}

// Create returns the element's syntax for a CREATE TABLE statement.
func (fk fkType) Create(d Dialect) (string, error) {
	// Compile the internal type, but only take the first word
	// TODO This is a hacky way to do this
	c, err := fk.t.Create(d)
	if err != nil {
		return "", err
	}
	t := strings.Split(c, " ")[0]
	return fmt.Sprintf(
		`%s REFERENCES %s("%s")`,
		t,
		fk.c.table.Name,
		fk.c.name,
	), nil
}

// ForeignKey creates a fkElem from the given name and column.
// The given column must already have an assigned table. The new ColumnElem
// will inherit its type from the given column's type.
// TODO Add the ability for self-referential foreign keys.
func ForeignKey(name string, fk ColumnElem) ColumnElem {
	return ColumnElem{
		inner: ColumnClause{name: name},
		name:  name,
		typ:   fkType{t: fk.typ, c: fk},
	}
}
