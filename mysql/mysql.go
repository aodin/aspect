package mysql

import (
	_ "github.com/go-sql-driver/mysql"

	"github.com/aodin/aspect"
)

// MySQL implements the Dialect interface for MySQL databases.
type MySQL struct{}

// Parameterize returns the MySQL specific parameterization scheme.
func (d *MySQL) Parameterize(i int) string {
	return `?`
}

// Add the mysql dialect to the dialect registry
func init() {
	aspect.RegisterDialect("mysql", &MySQL{})
}
