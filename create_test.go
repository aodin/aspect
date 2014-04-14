package aspect

import (
	"testing"
)

func TestCreateStmt(t *testing.T) {
	stmt := users.Create()
	expected := `CREATE TABLE "users" (
	"id" SERIAL NOT NULL,
	"name" VARCHAR(32) NOT NULL,
	"password" VARCHAR,
	PRIMARY KEY ("id")
)`
	expectedSQL(t, stmt, expected)
}
