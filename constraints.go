package aspect

import (
	"fmt"
	"strings"
)

// PrimaryKeyArray is a list of columns representing the table's primary key
// array. It implements the TableModifier and Creatable interfaces.
type PrimaryKeyArray []string

var _ Creatable = PrimaryKeyArray{}
var _ TableModifier = PrimaryKeyArray{}

// Create returns the proper syntax for CREATE TABLE commands.
func (pk PrimaryKeyArray) Create(d Dialect) (string, error) {
	cs := make([]string, len(pk))
	for i, c := range pk {
		cs[i] = fmt.Sprintf(`"%s"`, c)
	}
	return fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(cs, ", ")), nil
}

// Modify implements the TableModifier interface. It confirms that every column
// given exists in the parent table.
func (pk PrimaryKeyArray) Modify(table *TableElem) error {
	// Confirm that all columns in the primary key exists
	for _, name := range pk {
		_, exists := table.C[name]
		if !exists {
			return fmt.Errorf("No column with the name '%s' exists in the table '%s'. Is it declared after the PrimaryKey declaration?", name, table.Name)
		}
		// TODO Set the attributes of the column type?
	}
	// If all columns exist, set the primary key
	table.pk = pk

	// Add the pk to the create array
	table.creates = append(table.creates, pk)

	return nil
}

// Contains returns true if the PrimaryKeyArray contains the given column name.
func (pk PrimaryKeyArray) Contains(key string) bool {
	for _, name := range pk {
		if name == key {
			return true
		}
	}
	return false
}

// PrimaryKey creates a new PrimaryKeyArray. Only one primary key is allowed
// per table.
func PrimaryKey(names ...string) PrimaryKeyArray {
	return PrimaryKeyArray(names)
}

// UniqueConstraint is the internal representation of a UNIQUE constraint. It
// implements the TableModifier and Creatable interfaces.
type UniqueConstraint []string

var _ Creatable = UniqueConstraint{}
var _ TableModifier = UniqueConstraint{}

// Create returns the proper syntax for CREATE TABLE commands.
func (uc UniqueConstraint) Create(d Dialect) (string, error) {
	cs := make([]string, len(uc))
	for i, c := range uc {
		cs[i] = fmt.Sprintf(`"%s"`, c)
	}
	return fmt.Sprintf("UNIQUE (%s)", strings.Join(cs, ", ")), nil
}

// Modify implements the TableModifier interface. It confirms that every column
// given exists in the parent table.
func (uc UniqueConstraint) Modify(table *TableElem) error {
	// Confirm that all columns in the primary key exists
	for _, name := range uc {
		_, exists := table.C[name]
		if !exists {
			return fmt.Errorf("No column with the name '%s' exists in the table '%s'. Is it declared after the PrimaryKey declaration?", name, table.Name)
		}
		// TODO Set the attributes of the column type?
	}

	// Add the unique clause to the table
	table.uniques = append(table.uniques, uc)

	// Add the constraint to the table
	table.creates = append(table.creates, uc)

	return nil
}

// Unique creates a new UniqueConstraint from the given column names.
func Unique(names ...string) UniqueConstraint {
	return UniqueConstraint(names)
}
