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
	return m.n < 0
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
	return nil
}

func newMockResult(columns ...string) *Result {
	scanner := &mockScanner{
		columns: columns,
		n:       2,
	}
	return &Result{
		rows: scanner,
	}
}

type simpleUser struct {
	ID   int64
	Name string
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
	var user []simpleUser
	assert.Nil(result.All(&user))

	// Too many fields is fine if they're tagged
	result = newMockResult("id", "name")
	var extras []extraUser
	assert.Nil(result.All(&extras))

	// Test improper types
	// Non-pointer
	result = newMockResult("id", "name")
	var nonPointer []simpleUser
	assert.NotNil(result.All(nonPointer))

	// Non-slice pointer
	result = newMockResult("id", "name")
	var ids int64
	assert.NotNil(result.All(&ids))

	// TODO No tags and unequal number of fields

}
