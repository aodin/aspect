package aspect

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNilSelect = errors.New("select received a non-existant table or column")
)

// Selectable is an interface that allows both tables and columns to be
// selected. It is implemented by TableElem and ColumnElem.
type Selectable interface {
	Selectable() []ColumnElem
}

// SelectStmt is an internal representation of a SELECT statement.
// TODO Use clauses to build the parts of statement
// TODO Should this struct be exported?
type SelectStmt struct {
	tables  []*TableElem
	columns []ColumnElem
	joins   []JoinStmt
	cond    Clause
	groupBy []ColumnElem
	order   []OrderedColumn
	limit   int
	offset  int
	err     error
}

// TableExists checks if a table already exists in the SELECT statement.
func (stmt SelectStmt) TableExists(name string) bool {
	for _, table := range stmt.tables {
		if table.Name == name {
			return true
		}
	}
	return false
}

// TODO These need errors now
func (stmt SelectStmt) CompileTables(d Dialect, params *Parameters) []string {
	names := make([]string, len(stmt.tables))
	for i, table := range stmt.tables {
		names[i] = table.Compile(d, params)
	}
	return names
}

// TODO These need errors now
func (stmt SelectStmt) CompileColumns(d Dialect, params *Parameters) []string {
	names := make([]string, len(stmt.columns))
	for i, c := range stmt.columns {
		names[i], _ = c.Compile(d, params)
	}
	return names
}

// String will return the string representation of the statement using a
// default dialect.
func (stmt SelectStmt) String() string {
	compiled, _ := stmt.Compile(&defaultDialect{}, Params())
	return compiled
}

// Compile will compile the SELECT statement according to the given dialect.
func (stmt SelectStmt) Compile(d Dialect, params *Parameters) (string, error) {
	if stmt.err != nil {
		return "", stmt.err
	}
	compiled := fmt.Sprintf(
		"SELECT %s FROM %s",
		strings.Join(stmt.CompileColumns(d, params), ", "),
		strings.Join(stmt.CompileTables(d, params), ", "),
	)
	if stmt.joins != nil && len(stmt.joins) > 0 {
		for _, join := range stmt.joins {
			jc, err := join.Compile(d, params)
			if err != nil {
				return "", err
			}
			compiled += jc
		}
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

// Join adds a JOIN ... ON to the SELECT statement. The second column parameter
// will be used as the join table and will be removed from the statement's
// FROM selection.
// TODO A smarter result for determining the JOIN table.
func (stmt SelectStmt) Join(pre, post ColumnElem) SelectStmt {
	// Get the table of the pre element
	preTable := pre.Table()

	// If the preTable is not in the current select statement, add it
	if !stmt.TableExists(preTable.Name) {
		stmt.tables = append(stmt.tables, preTable)
	}

	// Get the table of the post element
	postTable := post.Table()

	// Remove the post table from the list of selected tables
	for i, t := range stmt.tables {
		if t == postTable {
			stmt.tables = append(stmt.tables[:i], stmt.tables[i+1:]...)
			break
		}
	}

	// Add the join structure to the select
	if stmt.joins == nil {
		stmt.joins = make([]JoinStmt, 0)
	}
	stmt.joins = append(
		stmt.joins,
		JoinStmt{method: "JOIN", pre: pre, post: post, table: postTable},
	)
	return stmt
}

// Where adds a conditional clause to the SELECT statement. Only one WHERE
// is allowed per statement. Additional calls to Where will overwrite the
// existing WHERE clause.
func (stmt SelectStmt) Where(cond Clause) SelectStmt {
	stmt.cond = cond
	return stmt
}

// GroupBy adds a GROUP BY to the SELECT statement. Only one GROUP BY
// is allowed per statement. Additional calls to GroupBy will overwrite the
// existing GROUP BY clause.
func (stmt SelectStmt) GroupBy(cs ...ColumnElem) SelectStmt {
	groupBy := make([]ColumnElem, len(cs))
	// Since columns may be given without an ordering method, perform the
	// orderable conversion whether or not it is already ordered
	for i, column := range cs {
		groupBy[i] = column
	}
	stmt.groupBy = groupBy
	return stmt
}

// OrderBy adds an ORDER BY to the SELECT statement. Only one ORDER BY
// is allowed per statement. Additional calls to OrderBy will overwrite the
// existing ORDER BY clause.
func (stmt SelectStmt) OrderBy(params ...Orderable) SelectStmt {
	order := make([]OrderedColumn, len(params))
	// Since columns may be given without an ordering method, perform the
	// orderable conversion whether or not it is already ordered
	for i, column := range params {
		order[i] = column.Orderable()
	}
	stmt.order = order
	return stmt
}

// Limit adds a LIMIT to the SELECT statement. Only one LIMIT is allowed per // statement. Additional calls to Limit will overwrite the existing LIMIT
// clause.
func (stmt SelectStmt) Limit(limit int) SelectStmt {
	stmt.limit = limit
	return stmt
}

// Offset adds an OFFSET to the SELECT statement. Only one OFFSET is allowed
// per statement. Additional calls to Offset will overwrite the existing
// OFFSET clause.
func (stmt SelectStmt) Offset(offset int) SelectStmt {
	stmt.offset = offset
	return stmt
}

// Select generates a new SELECT statement from the given columns and tables.
func Select(selections ...Selectable) (stmt SelectStmt) {
	// Iterate through the selections and get the columns in the selection
	for _, selection := range selections {
		if selection == nil {
			stmt.err = ErrNilSelect
			return
		}
		// Each selection may return multiple columns (as with a table)
		columns := selection.Selectable()
		for _, column := range columns {
			// Adding a bad column will pass a zero-initialized ColumnElem and
			// since blank column names are invalid SQL we can reject them
			if column.Name() == "" {
				stmt.err = ErrNilSelect
				return
			}
			stmt.columns = append(stmt.columns, column)

			// Add the table to the stmt tables if it does not already exist
			if !stmt.TableExists(column.Table().Name) {
				stmt.tables = append(stmt.tables, column.Table())
			}
		}
	}
	return
}

// SelectExcept generates a new SELECT statement from the given table except
// for the given columns.
func SelectExcept(table *TableElem, exceptions ...ColumnElem) SelectStmt {
	// TODO This func should proxy to Select in case behavior changes
	stmt := SelectStmt{
		tables: []*TableElem{table},
	}
	columns := table.Selectable()
	// Remove the exceptions
	// TODO Some set operations would be nice here
	for _, exception := range exceptions {
		for i, column := range columns {
			// There should only be one column matching the exception per table
			if exception == column {
				columns = append(columns[:i], columns[i+1:]...)
				break
			}
		}
	}
	stmt.columns = columns
	return stmt
}
