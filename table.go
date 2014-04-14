package aspect

import (
	"fmt"
)

/*

Table
-----

Create using the `Table()` constructor, `TableElem`s represent a Table within the SQL spec.

The Table consists of any number of `TableModifier` interfaces, including `Column`, `ForeignKey`, `Constraint`, and `PrimaryKey` elements.

It will panic when an improper schema is created.

*/

// Allow the implementing element to modify its parent Table.
// In the case of a ColumnStruct, it will add the column after checking
// for a namespace collision.
type TableModifier interface {
	Modify(*TableElem) error
}

type TableElem struct {
	Name  string
	C     ColumnSet // TODO Should this be a pointer?
	order []string
	pk    PrimaryKeyArray
}

func (table *TableElem) String() string {
	return table.Name
}

// TODO Compile might not be the best name for this method, since it is
// not a target for compilation
func (table *TableElem) Compile(d Dialect, params *Parameters) string {
	return fmt.Sprintf(`"%s"`, table.Name)
}

// Get the table columns in proper order
func (table *TableElem) Columns() []*ColumnStruct {
	columns := make([]*ColumnStruct, len(table.order))
	for index, name := range table.order {
		columns[index] = table.C[name]
	}
	return columns
}

func (table *TableElem) Create() *CreateStmt {
	return &CreateStmt{table: table}
}

func (table *TableElem) Drop() *DropStmt {
	return &DropStmt{table: table}
}

// Alias for Select(table) that will select all columns in the table
func (table *TableElem) Select() *SelectStatement {
	return Select(table)
}

// Implement the sql.Selectable interface for building SELECT statements
func (table *TableElem) Selectable() []ColumnElement {
	columns := make([]ColumnElement, len(table.order))
	for index, name := range table.order {
		columns[index] = table.C[name]
	}
	return columns
}

// Constructor Method for an DELETE statement tied to this table
func (table *TableElem) Delete(args ...interface{}) *DeleteStatement {
	return Delete(table, args...)
}

func (table *TableElem) Insert(arg interface{}, args ...interface{}) *InsertStatement {
	return InsertTableValues(table, arg, args...)
}

// Constructor function
func Table(name string, elements ...TableModifier) *TableElem {
	// Create the table with its dynamic elements
	mapping := ColumnSet{}
	order := make([]string, 0)

	// TODO Name safety
	table := &TableElem{
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
