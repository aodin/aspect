package aspect

import (
	"fmt"
	"strings"
)

/*

Primary Key
-----------

Implements the `TableModifier` interface.

*/

// Simply a list of columns
type PrimaryKeyArray []string

func (pk PrimaryKeyArray) Create(d Dialect) (string, error) {
	cs := make([]string, len(pk))
	for i, c := range pk {
		cs[i] = fmt.Sprintf(`"%s"`, c)
	}
	return fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(cs, ", ")), nil
}

// To implement the `TableModifier` interface, the struct must
// have method Modify(). It does not need to modify its parent table.
func (pk PrimaryKeyArray) Modify(table *TableElem) error {
	// Confirm that all columns in the primary key exists
	for _, name := range pk {
		// TODO Aggregate errors
		_, exists := table.C[name]
		if !exists {
			return fmt.Errorf("No column with the name '%s' exists in the table '%s'. Is it declared after the PrimaryKey declaration?", name, table.Name)
		}
		// TODO Set the attributes of the column type?
	}
	// If all columns exist, set the primary key
	table.pk = pk
	return nil
}

// TODO Compilation and string output

func (pk PrimaryKeyArray) Contains(key string) bool {
	for _, name := range pk {
		if name == key {
			return true
		}
	}
	return false
}

// Constructor function for PrimaryKeyArray
func PrimaryKey(names ...string) PrimaryKeyArray {
	return PrimaryKeyArray(names)
}
