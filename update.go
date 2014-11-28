package aspect

import (
	"fmt"
)

// UpdateStmt is the internal representation of an SQL UPDATE statement.
type UpdateStmt struct {
	table  *TableElem
	values Values
	err    error
	cond   Clause
}

func (stmt UpdateStmt) Error() error {
	return stmt.err
}

// String outputs the parameter-less UPDATE statement in a neutral dialect.
func (stmt UpdateStmt) String() string {
	compiled, _ := stmt.Compile(&defaultDialect{}, Params())
	return compiled
}

// Compile outputs the UPDATE statement using the given dialect and parameters.
// An error may be returned because of a pre-existing error or because
// an error occurred during compilation.
func (stmt UpdateStmt) Compile(d Dialect, params *Parameters) (string, error) {
	// Check for delayed errors
	if stmt.err != nil {
		return "", stmt.err
	}

	// If no values were attached, then create a default values map
	if stmt.values == nil {
		stmt.values = Values{}
		for _, column := range stmt.table.Columns() {
			stmt.values[column.Name()] = nil
		}
	}

	// Compile the values
	valuesStmt, err := stmt.values.Compile(d, params)
	if err != nil {
		return "", err
	}

	// Begin building the UPDATE statement
	compiled := fmt.Sprintf(
		`UPDATE "%s" SET %s`,
		stmt.table.Name,
		valuesStmt,
	)

	// Add a conditional statement if it exists
	if stmt.cond != nil {
		cc, err := stmt.cond.Compile(d, params)
		if err != nil {
			return "", err
		}
		compiled += fmt.Sprintf(" WHERE %s", cc)
	}
	return compiled, nil
}

// Values attaches the given values to the statement. The keys of values
// must match columns in the table.
// TODO Allow structs to be used if a primary key is specified in the schema
func (stmt UpdateStmt) Values(values Values) UpdateStmt {
	// There must be some columns to update!
	if len(values) == 0 {
		stmt.err = fmt.Errorf("aspect: there must be at least one value in a Values instance")
		return stmt
	}

	// Confirm that all values' keys are columns in the table
	for key := range values {
		if _, ok := stmt.table.C[key]; !ok {
			stmt.err = fmt.Errorf(
				"aspect: no column %s exists in the table %s",
				key,
				stmt.table.Name,
			)
		}
	}

	stmt.values = values
	return stmt
}

// Where adds a conditional WHERE clause to the UPDATE statement.
func (stmt UpdateStmt) Where(conds ...Clause) UpdateStmt {
	if len(conds) > 1 {
		// By default, multiple where clauses will be joined will AllOf
		stmt.cond = AllOf(conds...)
	} else if len(conds) == 1 {
		stmt.cond = conds[0]
	}
	return stmt
}

// Update creates an UPDATE statement for the given table.
func Update(table *TableElem) (stmt UpdateStmt) {
	if table == nil {
		stmt.err = fmt.Errorf(
			"aspect: attempting to UPDATE a nil table",
		)
		return
	}
	stmt.table = table
	return
}
