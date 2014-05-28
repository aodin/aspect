package aspect

import (
	"database/sql"
)

// TODO How to distiguish between full statements and fragments?
type Executable interface {
	Compiles
}

// TODO dialect
type DB struct {
	conn    *sql.DB
	dialect Dialect
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) Dialect() Dialect {
	return db.dialect
}

func (db *DB) Execute(stmt Executable, args ...interface{}) (*Result, error) {
	// Initialize a list of empty parameters
	params := Params()

	// TODO Columns are needed for name return types, tag matching, etc...
	s, err := stmt.Compile(db.dialect, params)
	if err != nil {
		return nil, err
	}

	// TODO When to use the given arguments?
	// TODO If args are structs, maps, or slices, unpack them
	// Use any arguments given in Execute() over compiled arguments
	if len(args) == 0 {
		args = params.args
	}

	rows, err := db.conn.Query(s, args...)
	if err != nil {
		return nil, err
	}
	// Wrap the sql rows in a result
	return &Result{rows: rows, stmt: s}, nil
}

// Execute the query statement and populate the interface with all results
func (db *DB) QueryAll(s Executable, i interface{}) error {
	result, err := db.Execute(s)
	if err != nil {
		return err
	}
	return result.All(i)
}

func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.conn.Exec(query, args...)
}

func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.conn.Query(query, args...)
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

func ConnectUsing(driver string, db *sql.DB) (*DB, error) {
	// Get the dialect
	dialect, err := GetDialect(driver)
	if err != nil {
		return nil, err
	}
	return &DB{conn: db, dialect: dialect}, nil
}

func Connect(driver, credentials string) (*DB, error) {
	// Connect to the database using the given credentials
	db, err := sql.Open(driver, credentials)
	if err != nil {
		return nil, err
	}

	// Get the dialect
	dialect, err := GetDialect(driver)
	if err != nil {
		return nil, err
	}
	return &DB{conn: db, dialect: dialect}, nil
}
