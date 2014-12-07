package postgres

import (
	"fmt"
	"github.com/aodin/aspect"
	_ "github.com/lib/pq"
)

// Implement the Dialect interface
type PostGres struct{}

func (d *PostGres) Parameterize(i int) string {
	return fmt.Sprintf(`$%d`, i)
}

// Add the postgres dialect to the dialect registry
func init() {
	aspect.RegisterDialect("postgres", &PostGres{})
}

// Now represents the clause needed to return a now timestamp in postgres
var Now string = `now() at time zone 'utc'`
