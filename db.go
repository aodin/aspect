package aspect

import (
	"database/sql"
)

// TODO The db should be able to determine if a stmt should be used with
// either Exec() or Query()

// TODO How to distiguish between full statements and fragments?
type Executable interface {
	Compiles
}

// DB includes the current sql.DB connection and the associated Dialect
type DB struct {
	conn    *sql.DB
	dialect Dialect
}

func (db *DB) Begin() (*TX, error) {
	tx, err := db.conn.Begin()
	return &TX{tx, db}, err
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) Dialect() Dialect {
	return db.dialect
}

func (db *DB) Query(stmt Executable, args ...interface{}) (*Result, error) {
	// Initialize a list of empty parameters
	params := Params()

	// TODO Columns are needed for name return types, tag matching, etc...
	s, err := stmt.Compile(db.dialect, params)
	if err != nil {
		return nil, err
	}

	// TODO When to use the given arguments?
	// TODO If args are structs, maps, or slices, unpack them
	// Use any arguments given to Query() over compiled arguments
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

// QueryAll will query the statement and populate the interface with all
// results
func (db *DB) QueryAll(s Executable, i interface{}) error {
	result, err := db.Query(s)
	if err != nil {
		return err
	}
	return result.All(i)
}

// QueryOne will query the statement and populate the interface with one result
func (db *DB) QueryOne(s Executable, i interface{}) error {
	result, err := db.Query(s)
	// Close the result rows or sqlite3 will open another connection
	defer result.rows.Close()
	if err != nil {
		return err
	}
	return result.One(i)
}

// Execute the statement
// Execute should be used when no results are expected
// TODO What should the return type be?
func (db *DB) Execute(stmt Executable, args ...interface{}) (sql.Result, error) {
	// Initialize a list of empty parameters
	params := Params()

	// TODO Columns are needed for name return types, tag matching, etc...
	s, err := stmt.Compile(db.dialect, params)
	if err != nil {
		return nil, err
	}

	// TODO When to use the given arguments?
	// TODO If args are structs, maps, or slices, unpack them
	// Use any arguments given to Query() over compiled arguments
	if len(args) == 0 {
		args = params.args
	}
	return db.conn.Exec(s, args...)
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

type TX struct {
	*sql.Tx
	db *DB
}

func (tx *TX) Query(stmt Executable, args ...interface{}) (*Result, error) {
	// Initialize a list of empty parameters
	params := Params()

	// TODO Columns are needed for name return types, tag matching, etc...
	s, err := stmt.Compile(tx.db.dialect, params)
	if err != nil {
		return nil, err
	}

	// TODO When to use the given arguments?
	// TODO If args are structs, maps, or slices, unpack them
	// Use any arguments given to Query() over compiled arguments
	if len(args) == 0 {
		args = params.args
	}

	rows, err := tx.Tx.Query(s, args...)
	if err != nil {
		return nil, err
	}
	// Wrap the sql rows in a result
	return &Result{rows: rows, stmt: s}, nil
}

// QueryAll will query the statement and populate the interface with all
// results
func (tx *TX) QueryAll(s Executable, i interface{}) error {
	result, err := tx.Query(s)
	if err != nil {
		return err
	}
	return result.All(i)
}

// QueryOne will query the statement and populate the interface with one result
func (tx *TX) QueryOne(s Executable, i interface{}) error {
	result, err := tx.Query(s)
	// Close the result rows or sqlite3 will open another connection
	defer result.rows.Close()
	if err != nil {
		return err
	}
	return result.One(i)
}

// Execute the statement
// Execute should be used when no results are expected
// TODO What should the return type be?
func (tx *TX) Execute(stmt Executable, args ...interface{}) (sql.Result, error) {
	// Initialize a list of empty parameters
	params := Params()

	// TODO Columns are needed for name return types, tag matching, etc...
	s, err := stmt.Compile(tx.db.dialect, params)
	if err != nil {
		return nil, err
	}

	// TODO When to use the given arguments?
	// TODO If args are structs, maps, or slices, unpack them
	// Use any arguments given to Query() over compiled arguments
	if len(args) == 0 {
		args = params.args
	}
	return tx.Exec(s, args...)
}
