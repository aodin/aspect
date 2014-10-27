package sqlite3

import (
	_ "github.com/mattn/go-sqlite3"

	"github.com/aodin/aspect"
)

// Sqlite3 implements the Dialect interface for sqlite3 databases.
type Sqlite3 struct{}

// Parameterize returns the sqlite3 specific parameterization scheme.
func (d *Sqlite3) Parameterize(i int) string {
	return `?`
}

// Add the sqlite3 dialect to the dialect registry
func init() {
	aspect.RegisterDialect("sqlite3", &Sqlite3{})
}
