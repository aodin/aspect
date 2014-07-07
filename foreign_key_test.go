package aspect

import (
	"testing"
)

var parents = Table("parents",
	Column("id", Integer{}),
	Column("name", String{}),
)

var children = Table("children",
	ForeignKey("parent_id", parents.C["id"]),
	Column("name", String{}),
)

func TestForeignKey(t *testing.T) {
	stmt := children.Create()
	expected := `CREATE TABLE "children" (
  "parent_id" INTEGER REFERENCES parents("id"),
  "name" VARCHAR
);`
	expectedSQL(t, stmt, expected, 0)
}
