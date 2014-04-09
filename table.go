package aspect

import (
	"fmt"
)

/*

Table
-----

Create using the `Table()` constructor, `TableStruct`s represent a Table within the SQL spec.

The Table consists of any number of `TableModifier` interfaces, including `Column`, `ForeignKey`, `Constraint`, and `PrimaryKey` elements.

It will panic when an improper schema is created.

*/

// Allow the implementing element to modify its parent Table.
// In the case of a ColumnStruct, it will add the column after checking
// for a namespace collision.
type TableModifier interface {
	Modify(*TableStruct) error
}

type TableStruct struct {
	Name  string
	C     ColumnSet // TODO Should this be a pointer?
	order []string
	pk    PrimaryKeyArray
}

func (table *TableStruct) String() string {
	return table.Name
}

func (table *TableStruct) Compile() string {
	return fmt.Sprintf(`"%s"`, table.Name)
}

// Alias for Select(table) that will select all columns in the table
func (table *TableStruct) Select() *SelectStatement {
	return Select(table)
}

// Implement the sql.Selectable interface for building SELECT statements
func (table *TableStruct) Selectable() []*ColumnStruct {
	// TODO This shouldn't have to be built everytime the table is selected
	columns := make([]*ColumnStruct, len(table.order))
	for index, name := range table.order {
		columns[index] = table.C[name]
	}
	return columns
}

// Constructor Method for an DELETE statement tied to this table
func (table *TableStruct) Delete() *DeleteStatement {
	return &DeleteStatement{Target: table}
}

// Implement the interface that is needed to generate a DELETE statement
func (table *TableStruct) Deletable() *TableStruct {
	return table
}

// Constructor function
func Table(name string, elements ...TableModifier) *TableStruct {
	// Create the table with its dynamic elements
	mapping := ColumnSet{}
	order := make([]string, 0)

	// TODO Name safety
	table := &TableStruct{
		Name:  name,
		C:     mapping,
		order: order,
	}

	// Pass the table to each element for potential modification
	for _, element := range elements {
		err := element.Modify(table)
		// Panic on error since little can be done with a bad schema
		// TODO Create an error-safe version of Table()? Use case?
		if err != nil {
			panic(err)
		}
	}

	return table
}
