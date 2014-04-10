package aspect

import (
	"fmt"
	"strings"
)

// Since an entire Table can be selected, this must return an array of columns
type Selectable interface {
	Selectable() []ColumnElement
}

type SelectStatement struct {
	tables  []*TableStruct
	columns []ColumnElement
	groupBy []ColumnElement
	order   []*OrderedColumn
	limit   int
	offset  int
}

// TODO Both the columns and tables should have their own methods for
// order, existance, and string compilation
// TODO Check the name or the physical column?
func (stmt *SelectStatement) ColumnExists(name string) bool {
	return false
}

func (stmt *SelectStatement) TableExists(name string) bool {
	for _, table := range stmt.tables {
		if table.Name == name {
			return true
		}
	}
	return false
}

func (stmt *SelectStatement) CompileTables() []string {
	names := make([]string, len(stmt.tables))
	for i, table := range stmt.tables {
		names[i] = table.Compile()
	}
	return names
}

func (stmt *SelectStatement) CompileColumns() []string {
	names := make([]string, len(stmt.columns))
	for i, c := range stmt.columns {
		names[i] = c.Compile()
	}
	return names
}

// TODO Use clauses to abstract compilation
func (stmt *SelectStatement) String() string {
	// TODO assume the default dialect if none is given
	// TODO Cache the result of compilation?
	return stmt.Compile()
}

// TODO Will require a dialect
func (stmt *SelectStatement) Compile() string {
	compiled := fmt.Sprintf(
		"SELECT %s FROM %s",
		strings.Join(stmt.CompileColumns(), ", "),
		strings.Join(stmt.CompileTables(), ", "),
	)
	if stmt.groupBy != nil && len(stmt.groupBy) > 0 {
		groupBy := make([]string, len(stmt.groupBy))
		for i, column := range stmt.groupBy {
			groupBy[i] = column.Compile()
		}
		compiled += fmt.Sprintf(" GROUP BY %s", strings.Join(groupBy, ", "))
	}
	if stmt.order != nil && len(stmt.order) > 0 {
		order := make([]string, len(stmt.order))
		for i, column := range stmt.order {
			order[i] = column.Compile()
		}
		compiled += fmt.Sprintf(" ORDER BY %s", strings.Join(order, ", "))
	}
	if stmt.limit != 0 {
		compiled += fmt.Sprintf(" LIMIT %d", stmt.limit)
	}
	if stmt.offset != 0 {
		compiled += fmt.Sprintf(" OFFSET %d", stmt.offset)
	}
	return compiled
}

func (stmt *SelectStatement) GroupBy(cs ...ColumnElement) *SelectStatement {
	groupBy := make([]ColumnElement, len(cs))
	// Since columns may be given without an ordering method, perform the
	// orderable conversion whether or not it is already ordered
	for i, column := range cs {
		groupBy[i] = column
	}
	stmt.groupBy = groupBy
	return stmt
}

func (stmt *SelectStatement) OrderBy(params ...Orderable) *SelectStatement {
	order := make([]*OrderedColumn, len(params))
	// Since columns may be given without an ordering method, perform the
	// orderable conversion whether or not it is already ordered
	for i, column := range params {
		order[i] = column.Orderable()
	}
	stmt.order = order
	return stmt
}

func (stmt *SelectStatement) Limit(limit int) *SelectStatement {
	stmt.limit = limit
	return stmt
}

func (stmt *SelectStatement) Offset(offset int) *SelectStatement {
	stmt.offset = offset
	return stmt
}

func (stmt *SelectStatement) Execute() (string, error) {
	// TODO Return any delayed errors
	// TODO Check for a cached string
	return stmt.Compile(), nil
}

func (stmt *SelectStatement) Args() []interface{} {
	return make([]interface{}, 0)
}

func Select(selections ...Selectable) *SelectStatement {
	stmt := &SelectStatement{
		columns: make([]ColumnElement, 0),
		tables:  make([]*TableStruct, 0),
	}

	// Iterate through the selections and get the columns in the selection
	for _, selection := range selections {
		columns := selection.Selectable()
		for _, column := range columns {
			// TODO Test for name conflicts
			stmt.columns = append(stmt.columns, column)

			if !stmt.TableExists(column.Table().Name) {
				stmt.tables = append(stmt.tables, column.Table())
			}
		}
	}
	return stmt
}
