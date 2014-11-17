package aspect

import (
	"testing"
)

func TestCreateStmt(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	expected := `CREATE TABLE "users" (
  "id" INTEGER NOT NULL,
  "name" VARCHAR(32) NOT NULL UNIQUE,
  "password" VARCHAR(128),
  PRIMARY KEY ("id")
);`
	expect.SQL(expected, users.Create())

	expected = `CREATE TABLE "attrs" (
  "id" INTEGER PRIMARY KEY,
  "a" INTEGER,
  "b" INTEGER,
  UNIQUE ("a", "b")
);`
	expect.SQL(expected, attrs.Create())
}
