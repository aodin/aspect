package postgres

import (
	"fmt"

	_ "github.com/lib/pq"

	"github.com/aodin/aspect"
)

// PostGres implements the Dialect interface for postgres databases.
type PostGres struct{}

// Parameterize returns the postgres specific parameterization scheme.
func (d *PostGres) Parameterize(i int) string {
	return fmt.Sprintf(`$%d`, i)
}

// Add the postgres dialect to the dialect registry
func init() {
	aspect.RegisterDialect("postgres", &PostGres{})
}

// Now represents the clause needed to return a now timestamp in postgres
var Now string = `now() at time zone 'utc'`
