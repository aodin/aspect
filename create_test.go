package aspect

import (
	"testing"
)

func TestCreateStmt(t *testing.T) {
	stmt := users.Create()
	expected := `CREATE TABLE "users" (
  "id" INTEGER NOT NULL,
  "name" VARCHAR(32) NOT NULL,
  "password" VARCHAR(128),
  PRIMARY KEY ("id")
);`
	expectedSQL(t, stmt, expected, 0)

	stmt = attrs.Create()
	expected = `CREATE TABLE "attrs" (
  "id" INTEGER PRIMARY KEY,
  "a" INTEGER,
  "b" INTEGER,
  UNIQUE ("a", "b")
);`
	expectedSQL(t, stmt, expected, 0)
}
