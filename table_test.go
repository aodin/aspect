package aspect

import (
	"testing"
)

// Declare schemas that can be used package-wide
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

var views = Table("views",
	Column("id", Integer{PrimaryKey: true}),
	Column("user_id", Integer{}),
	Column("url", String{}),
	Column("ip", Inet{}),
	Column("timestamp", Timestamp{}),
)

var edges = Table("edges",
	Column("a", Integer{}),
	Column("b", Integer{}),
	PrimaryKey("a", "b"),
)

type edge struct {
	A int64 `db:"a"`
	B int64 `db:"b"`
}

// A short test for testing that an SQL statement was compiled as expected
func expectedSQL(t *testing.T, stmt Compiles, expected string, p int) {
	params := Params()
	compiled, err := stmt.Compile(&defaultDialect{}, params)
	if err != nil {
		t.Error(err)
	}
	if compiled != expected {
		t.Errorf("Unexpected SQL: %s != %s", compiled, expected)
	}
	if params.Len() != p {
		t.Errorf(
			"Unexpected number of parameters for %s: %d != %d",
			expected,
			params.Len(),
			p,
		)
	}
}

func TestTableSchema(t *testing.T) {
	// Test table properties
	if users.Name != "users" {
		t.Errorf("Unexpected Table name: '%s' != 'users'", users.Name)
	}

	// Test the accessor methods
	userId := users.C["id"]
	if userId.name != "id" {
		t.Errorf("Name of returned column was not 'id': '%s'", userId.name)
	}

	// A pointer to the column's table should have been added
	if userId.table != users {
		t.Errorf("Unexpected Table: %+v", userId.table)
	}

	// TODO Test improper schemas
}

// Test the sql.SelectStmts generated by table.Select() and Select()
func TestTableSelect(t *testing.T) {
	// All three of these select statements should produce the same output
	x := `SELECT "users"."id", "users"."name", "users"."password" FROM "users"`
	expectedSQL(t, users.Select(), x, 0)
	expectedSQL(t, Select(users), x, 0)
	expectedSQL(
		t,
		Select(users.C["id"], users.C["name"], users.C["password"]),
		x,
		0,
	)
}

// Test the InsertStmt generated by table.Insert()
func TestTableInsert(t *testing.T) {
	// An example user
	admin := user{1, "admin", "secret"}

	// Insert a single value into the table
	stmt := users.Insert(&admin)
	expectedSQL(
		t,
		stmt,
		`INSERT INTO "users" ("id", "name", "password") VALUES ($1, $2, $3)`,
		3,
	)
}

// Test DeleteStmt generated by table.Delete()
func TestTableDelete(t *testing.T) {
	// Delete the entire table
	expectedSQL(t, users.Delete(), `DELETE FROM "users"`, 0)

	// Delete using a conditional
	expectedSQL(
		t,
		users.Delete().Where(users.C["id"].Equals(1)),
		`DELETE FROM "users" WHERE "users"."id" = $1`,
		1,
	)
}
