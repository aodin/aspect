package aspect

import "fmt"

// Stmt is the base of all statements, including SELECT, UPDATE, DELETE, and
// INSERT statements
type Stmt struct {
	err error
}

// Error returns the statement's inner error
func (stmt Stmt) Error() error {
	return stmt.err
}

func (stmt *Stmt) SetError(msg string, args ...interface{}) {
	stmt.err = fmt.Errorf(msg, args...)
}

// ConditionalStmt includes SELECT, DELETE, and UPDATE statements
type ConditionalStmt struct {
	Stmt
	cond Clause
}

// Conditional returns the statement's conditional Clause
func (stmt ConditionalStmt) Conditional() Clause {
	return stmt.cond
}

// AddConditional adds a conditional clause to the statement. If a conditional
// clause already exists, they will be joined with an AND.
func (stmt *ConditionalStmt) AddConditional(cond Clause) {
	if stmt.cond == nil {
		stmt.cond = cond
	} else {
		stmt.cond = AllOf(stmt.cond, cond)
	}
}
