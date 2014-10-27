package sqlite3

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"

	"github.com/aodin/aspect"
)

var users = aspect.Table("users",
	aspect.Column("id", aspect.Integer{NotNull: true}),
	aspect.Column("name", aspect.String{Length: 32, NotNull: true}),
	aspect.Column("password", aspect.String{Length: 128}),
	aspect.PrimaryKey("id"),
)

type user struct {
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	Password string `db:"password"`
}

// Connect to an in-memory sqlite3 instance and execute some statements.
func TestConnect(t *testing.T) {
	assert := assert.New(t)

	db, err := aspect.Connect("sqlite3", ":memory:")
	assert.Nil(err)
	defer db.Close()

	// Create the users table
	_, err = db.Execute(users.Create())
	assert.Nil(err)

	// Insert a user
	admin := user{
		ID:       1,
		Name:     "admin",
		Password: "secret",
	}
	_, err = db.Execute(users.Insert().Values(admin))
	assert.Nil(err)

	// Select a user - Query must be given a pointer
	var u user
	assert.Nil(db.QueryOne(users.Select(), &u))
	assert.Equal(admin.ID, u.ID)
	assert.Equal(admin.Name, u.Name)
	assert.Equal(admin.Password, u.Password)

	// Drop the table
	_, err = db.Execute(users.Drop())
	assert.Nil(err)
}
