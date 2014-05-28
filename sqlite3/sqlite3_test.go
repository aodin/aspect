package sqlite3

import (
	. "github.com/aodin/aspect"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

// TODO Explain the choice of sqlite3 driver

var users = Table("users",
	Column("id", Integer{NotNull: true}),
	Column("name", String{Length: 32, NotNull: true}),
	Column("password", String{Length: 128}),
	PrimaryKey("id"),
)

type user struct {
	Id       int64  `db:"id"`
	Name     string `db:"name"`
	Password string `db:"password"`
}

// Connect to an in-memory sqlite3 instance and execute some statements.
func TestConnect(t *testing.T) {
	db, err := Connect("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Create the users table
	_, err = db.Execute(users.Create())
	if err != nil {
		t.Fatal(err)
	}

	// Insert a user
	// Structs can be inserted by value or reference
	admin := user{Id: 1, Name: "admin", Password: "secret"}
	_, err = db.Execute(users.Insert(admin))
	if err != nil {
		t.Fatal(err)
	}

	// Select a user
	// Query must be given a pointer
	var u user
	err = db.QueryOne(users.Select(), &u)
	if err != nil {
		t.Fatal(err)
	}
	if u.Id != admin.Id {
		t.Errorf("Unexpected user id: %d", u.Id)
	}
	if u.Name != admin.Name {
		t.Errorf("Unexpected user name: %s", u.Name)
	}
	if u.Password != admin.Password {
		t.Errorf("Unexpected user password: %s", u.Password)
	}
}
