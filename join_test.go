package aspect

import (
	"testing"
)

var tableA = Table("a",
	Column("id", Integer{}),
	Column("value", String{}),
)

var tableB = Table("b",
	Column("a_id", Integer{}),
	Column("value", String{}),
)

func TestJoinStmt(t *testing.T) {
	a := Select(tableA, tableB).Join(tableA.C["id"], tableB.C["a_id"])
	expectedSQL(
		t,
		a,
		`SELECT "a"."id", "a"."value", "b"."a_id", "b"."value" FROM "a" JOIN "b" ON "a"."id" = "b"."a_id"`,
	)

	b := Select(tableA.C["value"], tableB.C["value"]).Join(tableA.C["id"], tableB.C["a_id"])
	expectedSQL(
		t,
		b,
		`SELECT "a"."value", "b"."value" FROM "a" JOIN "b" ON "a"."id" = "b"."a_id"`,
	)
}
