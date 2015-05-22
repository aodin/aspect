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
