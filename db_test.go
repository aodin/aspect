package aspect

import (
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

// Connect to an in-memory sqlite3 instance and execute some statements.
func TestConnect(t *testing.T) {
	conn, err := Connect("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// Create the users table
	// TODO Compile the statement and run on the sql Exec() method
	createUsers, err := CompileWith(users.Create(), conn.dialect)
	if err != nil {
		t.Fatal(err)
	}
	_, err = conn.Exec(createUsers)
	if err != nil {
		t.Fatal(err)
	}

	// Insert a user
	// TODO Compile the statement and run on the sql Exec() method
	admin := user{Id: 1, Name: "admin", Password: "secret"}
	params := Params()
	insertUser, err := CompileWithParams(
		users.Insert(admin),
		conn.dialect,
		params,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = conn.Exec(insertUser, params.args...)
	if err != nil {
		t.Fatal(err)
	}

	// Select a user
	var u user
	rows, err := conn.Execute(users.Select())
	if err != nil {
		t.Fatal(err)
	}
	err = rows.One(&u)
	if err != nil {
		t.Fatal(err)
	}
}
