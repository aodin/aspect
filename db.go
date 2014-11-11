package aspect

import (
	"database/sql"
)

// Connection is a common interface for database connections or transactions
type Connection interface {
	Begin() (Transaction, error)
	Execute(stmt Executable, args ...interface{}) (sql.Result, error)
	Query(stmt Executable, args ...interface{}) (*Result, error)
	QueryAll(stmt Executable, i interface{}) error
	QueryOne(stmt Executable, i interface{}) error
	String(stmt Executable) string // Parameter-less output for logging
}

// Both DB and TX should implement the Connection interface
var _ Connection = &DB{}
var _ Connection = &TX{}

type Transaction interface {
	Connection
	Commit() error
	Rollback() error
}

// Both TX and fakeTX should implement the Connection interface
var _ Transaction = &TX{}
var _ Transaction = &fakeTX{}

// TODO The db should be able to determine if a stmt should be used with
// either Exec() or Query()

// Executable statements implement the Compiles interface
type Executable interface {
	Compiles
}

// DB wraps the current sql.DB connection pool and includes the Dialect
// associated with the connection.
type DB struct {
	conn    *sql.DB
	dialect Dialect
}

// Begin starts a new transaction using the current database connection pool.
func (db *DB) Begin() (Transaction, error) {
	tx, err := db.conn.Begin()
	return &TX{Tx: tx, dialect: db.dialect}, err
}

// Close closes the current database connection pool.
func (db *DB) Close() error {
	return db.conn.Close()
}

// Dialect returns the dialect associated with the current database connection
// pool.
func (db *DB) Dialect() Dialect {
	return db.dialect
}

// Execute executes the Executable statement with optional arguments. It
// returns the database/sql package's Result object, which may contain
// information on rows affected and last ID inserted depending on the driver.
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

// Query executes an Executable statement with the optional arguments. It
// returns a Result object, that can scan rows in various data types.
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

// QueryAll will query the statement and populate the given interface with all
// results.
func (db *DB) QueryAll(stmt Executable, i interface{}) error {
	result, err := db.Query(stmt)
	if err != nil {
		return err
	}
	return result.All(i)
}

// QueryOne will query the statement and populate the given interface with a
// single result.
func (db *DB) QueryOne(stmt Executable, i interface{}) error {
	result, err := db.Query(stmt)
	if err != nil {
		return err
	}
	// Close the result rows or sqlite3 will open another connection
	defer result.rows.Close()
	return result.One(i)
}

// String returns parameter-less SQL. If an error occurred during compilation,
// then an empty string will be returned.
func (db *DB) String(stmt Executable) string {
	compiled, _ := stmt.Compile(db.dialect, Params())
	return compiled
}

// Connect connects to the database using the given driver and credentials.
// It returns a database connection pool and an error if one occurred.
func Connect(driver, credentials string) (*DB, error) {
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

// TX wraps the current sql.Tx transaction and the Dialect associated with
// the transaction.
type TX struct {
	*sql.Tx
	dialect Dialect
}

// Begin returns the existing transaction. TODO Are nested transactions
// possible? And on what dialects?
func (tx *TX) Begin() (Transaction, error) {
	return tx, nil
}

// Commit calls the wrapped transactions Commit method.
func (tx *TX) Commit() error {
	return tx.Tx.Commit()
}

// Query executes an Executable statement with the optional arguments
// using the current transaction. It returns a Result object, that can scan
// rows in various data types.
func (tx *TX) Query(stmt Executable, args ...interface{}) (*Result, error) {
	// Initialize a list of empty parameters
	params := Params()

	// TODO Columns are needed for name return types, tag matching, etc...
	s, err := stmt.Compile(tx.dialect, params)
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

// QueryAll will query the statement using the current transaction and
// populate the given interface with all results.
func (tx *TX) QueryAll(stmt Executable, i interface{}) error {
	result, err := tx.Query(stmt)
	if err != nil {
		return err
	}
	return result.All(i)
}

// QueryOne will query the statement using the current transaction and
// populate the given interface with a single result.
func (tx *TX) QueryOne(stmt Executable, i interface{}) error {
	result, err := tx.Query(stmt)
	if err != nil {
		return err
	}
	// Close the result rows or sqlite3 will open another connection
	defer result.rows.Close()
	return result.One(i)
}

// Execute executes the Executable statement with optional arguments using
// the current transaction. It returns the database/sql package's Result
// object, which may contain information on rows affected and last ID inserted
// depending on the driver.
func (tx *TX) Execute(stmt Executable, args ...interface{}) (sql.Result, error) {
	// Initialize a list of empty parameters
	params := Params()

	// TODO Columns are needed for name return types, tag matching, etc...
	s, err := stmt.Compile(tx.dialect, params)
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

// Rollback calls the wrapped transactions Rollback method.
func (tx *TX) Rollback() error {
	return tx.Tx.Rollback()
}

// String returns parameter-less SQL. If an error occurred during compilation,
// then an empty string will be returned.
func (tx *TX) String(stmt Executable) string {
	compiled, _ := stmt.Compile(tx.dialect, Params())
	return compiled
}

// WrapTx allows aspect to take control of an existing database/sql
// transaction and execute queries using the given dialect.
func WrapTx(tx *sql.Tx, dialect Dialect) *TX {
	return &TX{Tx: tx, dialect: dialect}
}

type fakeTX struct {
	tx Transaction
}

func (tx *fakeTX) Begin() (Transaction, error) {
	return tx, nil
}

func (tx *fakeTX) Commit() error {
	return nil
}

func (tx *fakeTX) Execute(stmt Executable, args ...interface{}) (sql.Result, error) {
	return tx.tx.Execute(stmt, args...)
}

func (tx *fakeTX) Query(stmt Executable, args ...interface{}) (*Result, error) {
	return tx.tx.Query(stmt, args...)
}

func (tx *fakeTX) QueryAll(stmt Executable, i interface{}) error {
	return tx.tx.QueryOne(stmt, i)
}

func (tx *fakeTX) QueryOne(stmt Executable, i interface{}) error {
	return tx.tx.QueryOne(stmt, i)
}

func (tx *fakeTX) Rollback() error {
	return nil
}

func (tx *fakeTX) String(stmt Executable) string {
	return tx.String(stmt)
}

// FakeTx allows testing of transactional blocks of code. Commit and Rollback
// do nothing.
func FakeTx(tx Transaction) *fakeTX {
	return &fakeTX{tx: tx}
}
