package aspect

import (
	"database/sql"
	"log"
)

// Connection is a common interface for database connections or transactions
type Connection interface {
	Begin() (Transaction, error)
	Execute(stmt Executable, args ...interface{}) (sql.Result, error)
	Query(stmt Executable, args ...interface{}) (*Result, error)
	QueryAll(stmt Executable, i interface{}) error
	QueryOne(stmt Executable, i interface{}) error
	String(stmt Executable) string // Parameter-less output for logging

	// Must operations will panic on error
	MustBegin() Transaction
	MustExecute(stmt Executable, args ...interface{}) sql.Result
	MustQuery(stmt Executable, args ...interface{}) *Result
	MustQueryAll(stmt Executable, i interface{})
	MustQueryOne(stmt Executable, i interface{}) bool
}

// Both DB and TX should implement the Connection interface
var _ Connection = &DB{}
var _ Connection = &TX{}

type Transaction interface {
	Connection
	Commit() error
	CommitIf(*bool) error
	MustCommitIf(*bool) bool
	Rollback() error
	MustRollbackIf(*bool)
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

func (db *DB) compile(stmt Executable) (string, *Parameters, error) {
	// Initialize a list of empty parameters
	params := Params()

	// Compile with the database connection's current dialect
	s, err := stmt.Compile(db.dialect, params)
	return s, params, err
}

// Execute executes the Executable statement with optional arguments. It
// returns the database/sql package's Result object, which may contain
// information on rows affected and last ID inserted depending on the driver.
func (db *DB) Execute(stmt Executable, args ...interface{}) (sql.Result, error) {
	s, params, err := db.compile(stmt)
	if err != nil {
		return nil, err
	}

	// Use any arguments given to Execute() over compiled arguments
	if len(args) == 0 {
		args = params.args
	}
	return db.conn.Exec(s, args...)
}

// Query executes an Executable statement with the optional arguments. It
// returns a Result object, that can scan rows in various data types.
func (db *DB) Query(stmt Executable, args ...interface{}) (*Result, error) {
	s, params, err := db.compile(stmt)
	if err != nil {
		return nil, err
	}

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
// then the string output of the error will be returned.
func (db *DB) String(stmt Executable) string {
	compiled, err := stmt.Compile(db.dialect, Params())
	if err != nil {
		return err.Error()
	}
	return compiled
}

// MustBegin starts a new transaction using the current database connection
// pool. It will panic on error.
func (db *DB) MustBegin() Transaction {
	tx, err := db.conn.Begin()
	if err != nil {
		log.Panicf(
			"aspect: failed to begin transaction: %s",
			err,
		)
	}
	return &TX{Tx: tx, dialect: db.dialect}
}

// MustExecute will panic on error.
func (db *DB) MustExecute(stmt Executable, args ...interface{}) sql.Result {
	s, params, err := db.compile(stmt)
	if err != nil {
		log.Panicf(
			"aspect: failed to compile (%s): %s",
			stmt,
			err,
		)
	}
	if len(args) == 0 {
		args = params.args
	}
	result, err := db.conn.Exec(s, args...)
	if err != nil {
		log.Panicf(
			"aspect: failed to exec (%s) with parameters (%v): %s",
			s,
			args,
			err,
		)
	}
	return result
}

// MustQuery
func (db *DB) MustQuery(stmt Executable, args ...interface{}) *Result {
	s, params, err := db.compile(stmt)
	if err != nil {
		log.Panicf(
			"aspect: failed to compile (%s): %s",
			stmt,
			err,
		)
	}

	// Use any arguments given to Query() over compiled arguments
	if len(args) == 0 {
		args = params.args
	}

	rows, err := db.conn.Query(s, args...)
	if err != nil {
		log.Panicf(
			"aspect: failed to query (%s) with parameters (%v): %s",
			s,
			args,
			err,
		)
	}
	// Wrap the sql rows in a result
	return &Result{rows: rows, stmt: s}
}

// MustQueryAll
func (db *DB) MustQueryAll(stmt Executable, i interface{}) {
	result := db.MustQuery(stmt)
	if err := result.All(i); err != nil {
		// TODO get parameters from result?
		log.Panicf(
			"aspect: failed to query all (%s): %s",
			result.stmt,
			err,
		)
	}
}

// MustQueryOne
func (db *DB) MustQueryOne(stmt Executable, i interface{}) bool {
	result := db.MustQuery(stmt)

	// Close the result rows or sqlite3 will open another connection
	defer result.rows.Close()

	err := result.One(i)
	if err == ErrNoResult {
		return false
	} else if err != nil {
		// TODO get parameters from result?
		log.Panicf(
			"aspect: failed to query one (%s): %s",
			result.stmt,
			err,
		)
	}
	return true
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

func (tx *TX) compile(stmt Executable) (string, *Parameters, error) {
	// Initialize a list of empty parameters
	params := Params()

	// Compile with the database connection's current dialect
	s, err := stmt.Compile(tx.dialect, params)
	return s, params, err
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

func (tx *TX) CommitIf(commit *bool) error {
	if *commit {
		return tx.Tx.Commit()
	} else {
		return tx.Tx.Rollback()
	}
}

func (tx *TX) MustCommitIf(commit *bool) bool {
	if *commit {
		if err := tx.Tx.Commit(); err != nil {
			log.Panicf("aspect: error during commit: %s", err)
		}
	} else {
		if err := tx.Tx.Rollback(); err != nil {
			log.Panicf("aspect: error during rollback: %s", err)
		}
	}
	return *commit
}

func (tx *TX) MustRollbackIf(rollback *bool) {
	if *rollback {
		if err := tx.Tx.Rollback(); err != nil {
			log.Panicf("aspect: error during rollback: %s", err)
		}
	}
}

// Execute executes the Executable statement with optional arguments using
// the current transaction. It returns the database/sql package's Result
// object, which may contain information on rows affected and last ID inserted
// depending on the driver.
func (tx *TX) Execute(stmt Executable, args ...interface{}) (sql.Result, error) {
	s, params, err := tx.compile(stmt)
	if err != nil {
		return nil, err
	}

	// Use any arguments given to Query() over compiled arguments
	if len(args) == 0 {
		args = params.args
	}
	return tx.Exec(s, args...)
}

// Query executes an Executable statement with the optional arguments
// using the current transaction. It returns a Result object, that can scan
// rows in various data types.
func (tx *TX) Query(stmt Executable, args ...interface{}) (*Result, error) {
	s, params, err := tx.compile(stmt)
	if err != nil {
		return nil, err
	}

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

// Begin returns the existing transaction.
func (tx *TX) MustBegin() Transaction {
	return tx
}

// MustExecute will panic on error.
func (tx *TX) MustExecute(stmt Executable, args ...interface{}) sql.Result {
	s, params, err := tx.compile(stmt)
	if err != nil {
		log.Panicf(
			"aspect: failed to compile (%s): %s",
			stmt,
			err,
		)
	}
	if len(args) == 0 {
		args = params.args
	}
	result, err := tx.Exec(s, args...)
	if err != nil {
		log.Panicf(
			"aspect: failed to exec (%s) with parameters (%v): %s",
			s,
			args,
			err,
		)
	}
	return result
}

// MustQuery
func (tx *TX) MustQuery(stmt Executable, args ...interface{}) *Result {
	s, params, err := tx.compile(stmt)
	if err != nil {
		log.Panicf(
			"aspect: failed to compile (%s): %s",
			stmt,
			err,
		)
	}

	// Use any arguments given to Query() over compiled arguments
	if len(args) == 0 {
		args = params.args
	}

	rows, err := tx.Tx.Query(s, args...)
	if err != nil {
		log.Panicf(
			"aspect: failed to query (%s) with parameters (%v): %s",
			s,
			args,
			err,
		)
	}
	// Wrap the sql rows in a result
	return &Result{rows: rows, stmt: s}
}

// MustQueryAll
func (tx *TX) MustQueryAll(stmt Executable, i interface{}) {
	result := tx.MustQuery(stmt)
	if err := result.All(i); err != nil {
		// TODO get parameters from result?
		log.Panicf(
			"aspect: failed to query all (%s): %s",
			result.stmt,
			err,
		)
	}
}

// MustQueryOne
func (tx *TX) MustQueryOne(stmt Executable, i interface{}) bool {
	result := tx.MustQuery(stmt)

	// Close the result rows or sqlite3 will open another connection
	defer result.rows.Close()

	err := result.One(i)
	if err == ErrNoResult {
		return false
	} else if err != nil {
		// TODO get parameters from result?
		log.Panicf(
			"aspect: failed to query one (%s): %s",
			result.stmt,
			err,
		)
	}
	return true
}

// Rollback calls the wrapped transactions Rollback method.
func (tx *TX) Rollback() error {
	return tx.Tx.Rollback()
}

// String returns parameter-less SQL. If an error occurred during compilation,
// then the string output of the error will be returned.
func (tx *TX) String(stmt Executable) string {
	compiled, err := stmt.Compile(tx.dialect, Params())
	if err != nil {
		return err.Error()
	}
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

func (tx *fakeTX) MustBegin() Transaction {
	return tx
}

func (tx *fakeTX) Begin() (Transaction, error) {
	return tx, nil
}

func (tx *fakeTX) Commit() error {
	return nil
}

func (tx *fakeTX) CommitIf(commit *bool) error {
	return nil
}

func (tx *fakeTX) MustCommitIf(commit *bool) bool {
	return false
}

func (tx *fakeTX) MustRollbackIf(rollback *bool) {}

func (tx *fakeTX) Execute(stmt Executable, args ...interface{}) (sql.Result, error) {
	return tx.tx.Execute(stmt, args...)
}

func (tx *fakeTX) Query(stmt Executable, args ...interface{}) (*Result, error) {
	return tx.tx.Query(stmt, args...)
}

func (tx *fakeTX) QueryAll(stmt Executable, i interface{}) error {
	return tx.tx.QueryAll(stmt, i)
}

func (tx *fakeTX) QueryOne(stmt Executable, i interface{}) error {
	return tx.tx.QueryOne(stmt, i)
}

func (tx *fakeTX) MustExecute(stmt Executable, args ...interface{}) sql.Result {
	return tx.tx.MustExecute(stmt, args...)
}

func (tx *fakeTX) MustQuery(stmt Executable, args ...interface{}) *Result {
	return tx.tx.MustQuery(stmt, args...)
}

func (tx *fakeTX) MustQueryAll(stmt Executable, i interface{}) {
	tx.tx.MustQueryAll(stmt, i)
}

func (tx *fakeTX) MustQueryOne(stmt Executable, i interface{}) bool {
	return tx.tx.MustQueryOne(stmt, i)
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
