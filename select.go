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
	cond    Clause
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

// TODO These need errors now
func (stmt *SelectStatement) CompileTables(d Dialect, params *Parameters) []string {
	names := make([]string, len(stmt.tables))
	for i, table := range stmt.tables {
		names[i] = table.Compile(d, params)
	}
	return names
}

// TODO These need errors now
func (stmt *SelectStatement) CompileColumns(d Dialect, params *Parameters) []string {
	names := make([]string, len(stmt.columns))
	for i, c := range stmt.columns {
		names[i], _ = c.Compile(d, params)
	}
	return names
}

func (stmt *SelectStatement) String() string {
	compiled, _ := stmt.Compile(&PostGres{}, Params())
	return compiled
}

// TODO Will require a dialect
func (stmt *SelectStatement) Compile(d Dialect, params *Parameters) (string, error) {
	compiled := fmt.Sprintf(
		"SELECT %s FROM %s",
		strings.Join(stmt.CompileColumns(d, params), ", "),
		strings.Join(stmt.CompileTables(d, params), ", "),
	)
	if stmt.cond != nil {
		cc, err := stmt.cond.Compile(d, params)
		if err != nil {
			return "", err
		}
		compiled += fmt.Sprintf(" WHERE %s", cc)
	}
	if stmt.groupBy != nil && len(stmt.groupBy) > 0 {
		groupBy := make([]string, len(stmt.groupBy))
		for i, column := range stmt.groupBy {
			// TODO Errors!
			groupBy[i], _ = column.Compile(d, params)
		}
		compiled += fmt.Sprintf(" GROUP BY %s", strings.Join(groupBy, ", "))
	}
	if stmt.order != nil && len(stmt.order) > 0 {
		order := make([]string, len(stmt.order))
		for i, column := range stmt.order {
			// TODO Errors!
			order[i], _ = column.Compile(d, params)
		}
		compiled += fmt.Sprintf(" ORDER BY %s", strings.Join(order, ", "))
	}
	if stmt.limit != 0 {
		compiled += fmt.Sprintf(" LIMIT %d", stmt.limit)
	}
	if stmt.offset != 0 {
		compiled += fmt.Sprintf(" OFFSET %d", stmt.offset)
	}
	return compiled, nil
}

func (stmt *SelectStatement) Where(cond Clause) *SelectStatement {
	stmt.cond = cond
	return stmt
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
