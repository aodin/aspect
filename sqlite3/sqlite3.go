package sqlite3

import (
	"github.com/aodin/aspect"
    _ "github.com/mattn/go-sqlite3"
)

// Implement the Dialect interface

type Sqlite3 struct{}

func (d *Sqlite3) Parameterize(i int) string {
	return `?`
}

// Add the sqlite3 dialect to the dialect registry
func init() {
	aspect.RegisterDialect("sqlite3", &Sqlite3{})
}
