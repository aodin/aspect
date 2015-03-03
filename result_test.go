package aspect

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO test not only length of args, but types
type mockScanner struct {
	columns []string
	n       int   // Default number of results
	err     error // Allow an error to be set
}

func (m *mockScanner) Close() error {
	return nil
}

func (m *mockScanner) Columns() ([]string, error) {
	return m.columns, nil
}

func (m *mockScanner) Err() error {
	return m.err
}

func (m *mockScanner) Next() bool {
	m.n--
	return m.n >= 0
}

func (m *mockScanner) Scan(args ...interface{}) error {
	// Arguments should match columns
	if len(m.columns) != len(args) {
		return fmt.Errorf(
			"aspect: expected %d arguments, got %d",
			len(m.columns),
			len(args),
		)
	}
	// Set the args
	// TODO Set the field to a non-empty value?
	for i, _ := range args {
		args[i] = nil
	}

	return nil
}

func newMockResult(columns ...string) *Result {
	return newMockResultN(2, columns...)
}

func newMockResultN(n int, columns ...string) *Result {
	scanner := &mockScanner{
		columns: columns,
		n:       n,
	}
	return &Result{
		rows: scanner,
	}
}

// Return an error on scan
type mockScanErrorResult struct {
	*mockScanner
}

func (m *mockScanErrorResult) Scan(args ...interface{}) error {
	return fmt.Errorf("aspect: mockScanErrorResult always errors")
}

func newMockScanErrorResult(columns ...string) *Result {
	scanner := &mockScanErrorResult{&mockScanner{
		columns: columns,
		n:       2,
	}}
	return &Result{
		rows: scanner,
	}
}

type simpleUser struct {
	ID   int64
	Name string
}

type tooManyFields struct {
	ID    int64
	Name  string
	Extra string
}

type extraUser struct {
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	Password string `db:"password"`
}

type outOfOrder struct {
	Password string `db:"password"`
	ID       int64  `db:"id"`
	Name     string `db:"name"`
}

func TestResult(t *testing.T) {
	assert := assert.New(t)

	// Test tagless unpacking
	result := newMockResult("id", "name")
	var simpleton simpleUser
	assert.Nil(result.One(&simpleton))

	result = newMockResult("id", "name")
	var simpletons []simpleUser
	assert.Nil(result.All(&simpletons))
	assert.Equal(2, len(simpletons))

	// Too many fields is fine if they're tagged
	result = newMockResult("id", "name")
	var extra extraUser
	assert.Nil(result.One(&extra))

	result = newMockResult("id", "name")
	var extras []extraUser
	assert.Nil(result.All(&extras))

	// Scan into arrays with elements
	result = newMockResult("id", "name")
	elems := []extraUser{{Password: "1"}, {Password: "2"}}
	assert.Nil(result.All(&elems))
	assert.Equal(2, len(elems))
	assert.Equal("1", elems[0].Password)

	// Test improper types
	// Non-pointer
	result = newMockResult("id", "name")
	var nonPointer []simpleUser
	assert.NotNil(result.All(nonPointer))

	// Non-slice pointer
	result = newMockResult("id", "name")
	var ids int64
	assert.NotNil(result.All(&ids))

	// No tags and unequal number of fields
	// 2015-03-03: fields branch - no longer an error
	// result = newMockResult("id", "name")
	// var tooManyOne tooManyFields
	// assert.NotNil(result.One(&tooManyOne))

	// result = newMockResult("id", "name")
	// var tooManyAll []tooManyFields
	// assert.NotNil(result.All(&tooManyAll))

	// Return no results
	result = newMockResultN(0, "id", "name")
	var none []tooManyFields
	assert.NotNil(result.One(&none))

	// Handle scan errors
	result = newMockScanErrorResult("id", "name")
	var oneError simpleUser
	assert.NotNil(result.One(&oneError))

	result = newMockScanErrorResult("id", "name")
	var allErrors []simpleUser
	assert.NotNil(result.All(&allErrors))

	// Attempt to scan into an array with One
	result = newMockScanErrorResult("id", "name")
	var simpletonErrors []simpleUser
	assert.NotNil(result.One(&simpletonErrors))
}
