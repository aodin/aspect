package postgres

import (
	"fmt"
	"strings"

	"github.com/aodin/aspect"
)

// RetInsertStmt is the internal representation of an INSERT ... RETURNING
// statement.
type RetInsertStmt struct {
	aspect.InsertStmt
	returning []aspect.ColumnElem
}

// String outputs the parameter-less INSERT ... RETURNING statement in a
// neutral dialect.
func (stmt RetInsertStmt) String() string {
	compiled, _ := stmt.Compile(&PostGres{}, aspect.Params())
	return compiled
}

// Compile outputs the INSERT ... RETURNING statement using the given dialect
// and parameters. An error may be returned because of a pre-existing error
// or because an error occurred during compilation.
func (stmt RetInsertStmt) Compile(d aspect.Dialect, params *aspect.Parameters) (string, error) {
	compiled, err := stmt.InsertStmt.Compile(d, params)
	if err != nil {
		return "", err
	}
	if len(stmt.returning) > 0 {
		compiled += fmt.Sprintf(
			" RETURNING %s",
			strings.Join(
				stmt.CompileColumns(d, params, stmt.returning),
				", ",
			),
		)
	}
	return compiled, nil
}

// TODO are errors required?
// TODO This should be a generalized function
func (stmt RetInsertStmt) CompileColumns(d aspect.Dialect, params *aspect.Parameters, columns []aspect.ColumnElem) []string {
	names := make([]string, len(columns))
	for i, c := range columns {
		names[i], _ = c.Compile(d, params)
	}
	return names
}

// Returning adds a RETURNING clause to the statement.
func (stmt RetInsertStmt) Returning(cs ...aspect.ColumnElem) RetInsertStmt {
	for _, column := range cs {
		if column.Table() != stmt.Table() {
			stmt.SetError(fmt.Errorf(
				"postgres: all columns must belong to the same table",
			))
			break
		}
		stmt.returning = append(stmt.returning, column)
	}
	return stmt
}

// Values proxies to the inner InsertStmt's Values method
func (stmt RetInsertStmt) Values(args interface{}) RetInsertStmt {
	stmt.InsertStmt = stmt.InsertStmt.Values(args)
	return stmt
}

// Insert creates an INSERT ... RETURNING statement for the given columns.
// There must be at least one column and all columns must belong to the
// same table.
func Insert(c aspect.ColumnElem, cs ...aspect.ColumnElem) RetInsertStmt {
	return RetInsertStmt{
		InsertStmt: aspect.Insert(c, cs...),
	}
}
