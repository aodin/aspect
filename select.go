package aspect

import (
	"fmt"
	"strings"
)

// Selectable is an interface that allows both tables and columns to be
// selected. It is implemented by TableElem and ColumnElem.
type Selectable interface {
	Selectable() []ColumnElem
}

// SelectStmt is the internal representation of an SQL SELECT statement.
type SelectStmt struct {
	tables  []*TableElem
	columns []ColumnElem
	joins   []JoinStmt // Broken and deprecated
	join    []JoinOnStmt
	cond    Clause
	groupBy []ColumnElem
	order   []OrderedColumn
	limit   int
	offset  int
	err     error // TODO common error handling struct
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

func (stmt SelectStmt) CompileTables(d Dialect, params *Parameters) []string {
	names := make([]string, len(stmt.tables))
	for i, table := range stmt.tables {
		names[i] = table.Compile(d, params)
	}
	return names
}

func (stmt SelectStmt) CompileColumns(d Dialect, params *Parameters) []string {
	names := make([]string, len(stmt.columns))
	for i, c := range stmt.columns {
		names[i], _ = c.Compile(d, params)
	}
	return names
}

// String outputs the parameter-less SELECT statement in a neutral dialect.
func (stmt SelectStmt) String() string {
	compiled, _ := stmt.Compile(&defaultDialect{}, Params())
	return compiled
}

// Compile outputs the SELECT statement using the given dialect and parameters.
// An error may be returned because of a pre-existing error or because
// an error occurred during compilation.
func (stmt SelectStmt) Compile(d Dialect, params *Parameters) (string, error) {
	if stmt.err != nil {
		return "", stmt.err
	}
	compiled := fmt.Sprintf(
		"SELECT %s FROM %s",
		strings.Join(stmt.CompileColumns(d, params), ", "),
		strings.Join(stmt.CompileTables(d, params), ", "),
	)

	// The old join syntax, which is broken and deprecated
	if len(stmt.joins) > 0 {
		for _, join := range stmt.joins {
			jc, err := join.Compile(d, params)
			if err != nil {
				return "", err
			}
			compiled += jc
		}
	}

	// JOIN ... ON ...
	if len(stmt.join) > 0 {
		for _, j := range stmt.join {
			jc, err := j.Compile(d, params)
			if err != nil {
				return "", err
			}
			compiled += jc
		}
	}

	// WHERE ...
	if stmt.cond != nil {
		cc, err := stmt.cond.Compile(d, params)
		if err != nil {
			return "", err
		}
		compiled += fmt.Sprintf(" WHERE %s", cc)
	}

	// GROUP BY ...
	if len(stmt.groupBy) > 0 {
		groupBy := make([]string, len(stmt.groupBy))
		for i, column := range stmt.groupBy {
			// TODO Errors!
			groupBy[i], _ = column.Compile(d, params)
		}
		compiled += fmt.Sprintf(" GROUP BY %s", strings.Join(groupBy, ", "))
	}

	// ORDER BY ...
	if len(stmt.order) > 0 {
		order := make([]string, len(stmt.order))
		for i, column := range stmt.order {
			// TODO Errors!
			order[i], _ = column.Compile(d, params)
		}
		compiled += fmt.Sprintf(" ORDER BY %s", strings.Join(order, ", "))
	}

	// LIMIT ...
	if stmt.limit != 0 {
		compiled += fmt.Sprintf(" LIMIT %d", stmt.limit)
	}

	// OFFSET ...
	if stmt.offset != 0 {
		compiled += fmt.Sprintf(" OFFSET %d", stmt.offset)
	}
	return compiled, nil
}

// Join adds a JOIN ... ON to the SELECT statement. The second column parameter
// will be used as the join table and will be removed from the statement's
// FROM selection.
// This method is broken and deprecated.
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
	stmt.joins = append(
		stmt.joins,
		JoinStmt{method: "JOIN", pre: pre, post: post, table: postTable},
	)
	return stmt
}

// Join adds a JOIN ... ON ... clause to the SELECT statement.
// TODO At least on clause should be required
func (stmt SelectStmt) JoinOn(table *TableElem, clauses ...Clause) SelectStmt {
	stmt.join = append(
		stmt.join,
		JoinOnStmt{
			method:      "JOIN",
			table:       table,
			ArrayClause: ArrayClause{Clauses: clauses, Sep: " AND "},
		},
	)
	return stmt
}

// Join adds a LEFT OUTER JOIN ... ON ... clause to the SELECT statement.
// TODO At least on clause should be required
func (stmt SelectStmt) LeftOuterJoinOn(table *TableElem, clauses ...Clause) SelectStmt {
	stmt.join = append(
		stmt.join,
		JoinOnStmt{
			method:      "LEFT OUTER JOIN",
			table:       table,
			ArrayClause: ArrayClause{Clauses: clauses, Sep: " AND "},
		},
	)
	return stmt
}

// Where adds a conditional clause to the SELECT statement. Only one WHERE
// is allowed per statement. Additional calls to Where will overwrite the
// existing WHERE clause.
func (stmt SelectStmt) Where(conds ...Clause) SelectStmt {
	if len(conds) > 1 {
		// By default, multiple where clauses will be joined will AllOf
		stmt.cond = AllOf(conds...)
	} else if len(conds) == 1 {
		stmt.cond = conds[0]
	}
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

// From allows the SelectStmt's FROM clause to be manually specified
// Since selections and joins will change the statement's currently selected
// tables, this method should be added to the end of a selection chain.
func (stmt SelectStmt) From(tables ...*TableElem) SelectStmt {
	stmt.tables = tables
	return stmt
}

// SelectTable creates a SELECT statement from the given table and its
// columns. Any additional selections will not have their table added to
// the SelectStmt's tables field - they must be added with the JoinOn syntax.
// To perform selections using cartesian logic, use Select() instead.
func SelectTable(table *TableElem, selects ...Selectable) (stmt SelectStmt) {
	stmt.tables = []*TableElem{table}
	stmt.columns = table.Columns()

	// Add any additional selections
	for _, selection := range selects {
		if selection == nil {
			stmt.err = fmt.Errorf("aspect: received a nil selectable - do the columns or tables you selected exist?")
			return
		}
		stmt.columns = append(stmt.columns, selection.Selectable()...)
	}

	// Confirm that all tables exist
	for _, column := range stmt.columns {
		if column.Name() == "" {
			stmt.err = fmt.Errorf("aspect: selected column does not exist")
			return
		}
	}
	return
}

// Select generates a new SELECT statement from the given columns and tables.
func Select(selections ...Selectable) (stmt SelectStmt) {
	columns := make([]ColumnElem, 0)
	for _, selection := range selections {
		if selection == nil {
			stmt.err = fmt.Errorf("aspect: received a nil selectable - do the columns or tables you selected exist?")
			return
		}
		columns = append(columns, selection.Selectable()...)
	}

	if len(columns) < 1 {
		stmt.err = fmt.Errorf("aspect: must select at least one column")
		return
	}

	for _, column := range columns {
		// Adding a bad column will pass a zero-initialized ColumnElem and
		// since blank column names are invalid SQL we can reject them
		if column.Name() == "" {
			stmt.err = fmt.Errorf("aspect: selected column does not exist")
			return
		}
		stmt.columns = append(stmt.columns, column)

		// Add the table to the stmt tables if it does not already exist
		if !stmt.TableExists(column.Table().Name) {
			stmt.tables = append(stmt.tables, column.Table())
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
