package aspect

import (
	"fmt"
)

// UpdateStmt represents an SQL UPDATE statement.
type UpdateStmt struct {
	table  *TableElem
	values Values
	err    error
	cond   Clause
}

// String returns an UPDATE statement using a default dialect.
// Parameters are discarded.
func (stmt UpdateStmt) String() string {
	compiled, _ := stmt.Compile(&defaultDialect{}, Params())
	return compiled
}

// Compile returns an UPDATE statement using the given dialect.
func (stmt UpdateStmt) Compile(d Dialect, params *Parameters) (string, error) {
	// Check for delayed errors
	if stmt.err != nil {
		return "", stmt.err
	}
	// Compile the values
	valuesStmt, err := stmt.values.Compile(d, params)
	if err != nil {
		return "", err
	}
	// Begin building the SQL UPDATE statement
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

// Where adds a conditional statement to the UPDATE statement.
func (stmt UpdateStmt) Where(cond Clause) UpdateStmt {
	stmt.cond = cond
	return stmt
}

// Update creates an UPDATE statement for the given table and values.
func Update(table *TableElem, values Values) (stmt UpdateStmt) {
	// TODO Or Update(TableElem, interface{}) for multiple or single structs
	stmt.table = table
	stmt.values = values

	// There must be some columns to update!
	if len(values) == 0 {
		stmt.err = ErrNoColumns
		return
	}

	// Confirm that all values' keys are columns in the table
	for key := range values {
		if _, ok := table.C[key]; !ok {
			stmt.err = fmt.Errorf(
				`aspect: no column "%s" exists in the table "%s"`,
				key,
				table.Name,
			)
		}
	}
	return stmt
}
