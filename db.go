package aspect

import (
	"database/sql"
)

type Executable interface {
	Execute() (string, error)
}

// TODO dialect
type DB struct {
	conn *sql.DB
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) Execute(stmt Executable, args ...interface{}) (*Result, error) {
	// TODO A dialect is needed to perform parameterization
	// TODO Columns are needed for name return types, tag matching, etc...
	s, err := stmt.Execute()
	if err != nil {
		return nil, err
	}

	// TODO If params are structs, maps, or slices, unpack them

	rows, err := db.conn.Query(s, args...)
	if err != nil {
		return nil, err
	}

	// Wrap the sql rows in a result
	return &Result{rows: rows, stmt: s}, nil
}

// A version of Execute that will panic if there is an error
func (db *DB) MustExecute(stmt Executable, args ...interface{}) *Result {
	result, err := db.Execute(stmt, args...)
	if err != nil {
		panic(err)
	}
	return result
}

func (db *DB) ExecuteSQL(s string, args ...interface{}) (*Result, error) {
	// TODO User Exec if no args are given?
	rows, err := db.conn.Query(s, args...)
	if err != nil {
		return nil, err
	}

	// Wrap the sql rows in a result
	return &Result{rows: rows, stmt: s}, nil
}

// A version of Execute that will panic if there is an error
func (db *DB) MustExecuteSQL(s string, args ...interface{}) *Result {
	result, err := db.ExecuteSQL(s, args...)
	if err != nil {
		panic(err)
	}
	return result
}

func Connect(driver, credentials string) (*DB, error) {
	db, err := sql.Open(driver, credentials)
	return &DB{conn: db}, err
}
