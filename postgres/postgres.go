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
