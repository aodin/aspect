package aspect

import (
	"fmt"
	"log"
)

// TableModifer is an interface with a single method Modify, which allows
// elements declared in a table schema to modify the parent table.
// In the case of a ColumnElem, it will add the column after checking
// for a namespace collision.
// TODO make this an internal interface?
type TableModifier interface {
	Modify(*TableElem) error
}

// Table names must validate
func validateTableName(name string) error {
	// TODO more rules
	if name == "" {
		return fmt.Errorf("aspect: table names cannot be blank")
	}
	return nil
}

// TableElem is underlying struct that holds an SQL TABLE schema. It is
// returned from the Table constructor function.
// TODO make this an internal struct?
type TableElem struct {
	Name    string
	C       ColumnSet
	order   []string
	pk      PrimaryKeyArray
	fks     []ForeignKeyElem
	uniques []UniqueConstraint
	creates []Creatable
}

// String returns the table name.
func (table *TableElem) String() string {
	return table.Name
}

// PrimaryKey returns the table's primary key array.
func (table *TableElem) PrimaryKey() PrimaryKeyArray {
	return table.pk
}

// UniqueConstraints returns the table's unique constraints.
func (table *TableElem) UniqueConstraints() []UniqueConstraint {
	return table.uniques
}

// ForeignKeys returns the table's foreign keys.
func (table *TableElem) ForeignKeys() []ForeignKeyElem {
	return table.fks
}

// Compile implements the Compiles interface allowing its use in statements.
// TODO Compile might not be the best name for this method, since it is
// not a target for compilation
func (table *TableElem) Compile(d Dialect, params *Parameters) string {
	return fmt.Sprintf(`"%s"`, table.Name)
}

// Columns returns the table's columns in proper order.
func (table *TableElem) Columns() []ColumnElem {
	columns := make([]ColumnElem, len(table.order))
	for index, name := range table.order {
		columns[index] = table.C[name]
	}
	return columns
}

// Create generates the table's CREATE statement.
func (table *TableElem) Create() CreateStmt {
	return CreateStmt{table: table}
}

// Create generates the table's DROP statement.
func (table *TableElem) Drop() DropStmt {
	return DropStmt{table: table}
}

// Select is an alias for Select(table). It will generate a SELECT statement
// for all columns in the table.
func (table *TableElem) Select(selections ...Selectable) SelectStmt {
	// If additional selections were provide, use the new behavior of
	// selection - this will not add the additional selections to the
	// SelectStmt's FROM clause
	if len(selections) > 0 {
		return SelectTable(table, selections...)
	}

	// Otherwise, use the old behavior
	return Select(table)
}

// SelectExcept is an alias for SelectExcept(table). It will generate a SELECT
// statement for all columns in the table except those provided as parameters.
func (table *TableElem) SelectExcept(exceptions ...ColumnElem) SelectStmt {
	return SelectExcept(table, exceptions...)
}

// Selectable allows the table to implement the Selectable interface, which
// builds SELECT statements.
func (table *TableElem) Selectable() []ColumnElem {
	columns := make([]ColumnElem, len(table.order))
	for index, name := range table.order {
		columns[index] = table.C[name]
	}
	return columns
}

// Delete is an alias for Delete(table). It will generate a DELETE statement
// for the entire table. Conditionals can be added with the Where() method or
// by specifying structs or slices of structs with Values()
func (table *TableElem) Delete() DeleteStmt {
	return Delete(table)
}

// Insert is an alias for Insert(table). It will create an INSERT statement
// for the entire table. Specify the insert values with the method Values().
func (table *TableElem) Insert() InsertStmt {
	return Insert(table)
}

// Update is an alias for Update(table). It will create an UPDATE statement
// for the entire table. Specify the update values with the method Values().
func (table *TableElem) Update() UpdateStmt {
	return Update(table)
}

// Table is the constructor function for TableElem. It is provided a name and
// any number of columns and constraints that implement the TableModifier
// interface.
func Table(name string, elements ...TableModifier) *TableElem {
	if err := validateTableName(name); err != nil {
		log.Panic(err)
	}

	table := &TableElem{
		Name: name,
		C:    ColumnSet{},
	}

	// Pass the table to each element for potential modification
	for _, element := range elements {
		// Panic on error since little can be done with a bad schema
		if err := element.Modify(table); err != nil {
			log.Panic(err)
		}
	}
	return table
}
