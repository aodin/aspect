package aspect

import (
	"fmt"
	"strings"
)

// Since an entire Table can be selected, this must return an array of columns
type Selectable interface {
	Selectable() []ColumnStruct
}

// TODO Use clauses to build the parts of statement
type SelectStmt struct {
	tables  []*TableElem
	columns []ColumnStruct
	join    *JoinStmt
	cond    Clause
	groupBy []ColumnStruct
	order   []OrderedColumn
	limit   int
	offset  int
}

// TODO Both the columns and tables should have their own methods for
// order, existance, and string compilation
// TODO Check the name or the physical column?
func (stmt *SelectStmt) ColumnExists(name string) bool {
	return false
}

func (stmt *SelectStmt) TableExists(name string) bool {
	for _, table := range stmt.tables {
		if table.Name == name {
			return true
		}
	}
	return false
}

// TODO These need errors now
func (stmt *SelectStmt) CompileTables(d Dialect, params *Parameters) []string {
	names := make([]string, len(stmt.tables))
	for i, table := range stmt.tables {
		names[i] = table.Compile(d, params)
	}
	return names
}

// TODO These need errors now
func (stmt *SelectStmt) CompileColumns(d Dialect, params *Parameters) []string {
	names := make([]string, len(stmt.columns))
	for i, c := range stmt.columns {
		names[i], _ = c.Compile(d, params)
	}
	return names
}

func (stmt *SelectStmt) String() string {
	compiled, _ := stmt.Compile(&PostGres{}, Params())
	return compiled
}

func (stmt *SelectStmt) Compile(d Dialect, params *Parameters) (string, error) {
	compiled := fmt.Sprintf(
		"SELECT %s FROM %s",
		strings.Join(stmt.CompileColumns(d, params), ", "),
		strings.Join(stmt.CompileTables(d, params), ", "),
	)
	if stmt.join != nil {
		jc, err := stmt.join.Compile(d, params)
		if err != nil {
			return "", err
		}
		compiled += jc
	}
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

// Add a JOIN ... ON to the SELECT statement
func (stmt *SelectStmt) Join(pre, post ColumnStruct) *SelectStmt {
	// Get the table of the post element
	table := post.Table()

	// Remove the post table from the list of selected tables
	for i, t := range stmt.tables {
		if t == table {
			stmt.tables = append(stmt.tables[:i], stmt.tables[i+1:]...)
			break
		}
	}

	// Add the join structure to the select
	stmt.join = &JoinStmt{pre: pre, post: post, table: table}
	return stmt
}

func (stmt *SelectStmt) Where(cond Clause) *SelectStmt {
	stmt.cond = cond
	return stmt
}

// Add a GROUP BY to the SELECT statement
func (stmt *SelectStmt) GroupBy(cs ...ColumnStruct) *SelectStmt {
	groupBy := make([]ColumnStruct, len(cs))
	// Since columns may be given without an ordering method, perform the
	// orderable conversion whether or not it is already ordered
	for i, column := range cs {
		groupBy[i] = column
	}
	stmt.groupBy = groupBy
	return stmt
}

// Add an ORDER BY to the SELECT statement
func (stmt *SelectStmt) OrderBy(params ...Orderable) *SelectStmt {
	order := make([]OrderedColumn, len(params))
	// Since columns may be given without an ordering method, perform the
	// orderable conversion whether or not it is already ordered
	for i, column := range params {
		order[i] = column.Orderable()
	}
	stmt.order = order
	return stmt
}

// Add a LIMIT to the SELECT statement
func (stmt *SelectStmt) Limit(limit int) *SelectStmt {
	stmt.limit = limit
	return stmt
}

// Add an OFFSET to the SELECT statement
func (stmt *SelectStmt) Offset(offset int) *SelectStmt {
	stmt.offset = offset
	return stmt
}

func Select(selections ...Selectable) *SelectStmt {
	stmt := &SelectStmt{
		columns: make([]ColumnStruct, 0),
		tables:  make([]*TableElem, 0),
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
