package aspect

import (
	"testing"
)

// TODO All table schemas should live in one file

var parents = Table("parents",
	Column("id", Integer{}),
	Column("name", String{}),
)

var children = Table("children",
	ForeignKey("parent_id", parents.C["id"]),
	Column("name", String{}),
)

func TestForeignKey_Create(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	expected := `CREATE TABLE "children" (
  "parent_id" INTEGER REFERENCES parents("id"),
  "name" VARCHAR
);`
	expect.SQL(expected, children.Create())
}
