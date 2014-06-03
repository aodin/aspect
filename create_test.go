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
}
