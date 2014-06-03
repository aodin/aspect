package aspect

import (
	"testing"
)

var tableA = Table("a",
	Column("id", Integer{}),
	Column("value", String{}),
)

var tableB = Table("b",
	Column("id", Integer{}),
	Column("value", String{}),
)

var relations = Table("relations",
	Column("a_id", Integer{}),
	Column("b_id", Integer{}),
	Unique("a_id", "b_id"),
)

func TestJoinStmt(t *testing.T) {
	a := Select(tableA).Join(tableA.C["id"], relations.C["a_id"])
	expectedSQL(
		t,
		a,
		`SELECT "a"."id", "a"."value" FROM "a" JOIN "relations" ON "a"."id" = "relations"."a_id"`,
		0,
	)

	b := Select(
		tableA.C["value"],
		tableB.C["value"],
	).Join(
		relations.C["a_id"],
		tableA.C["id"],
	).Join(
		relations.C["b_id"],
		tableB.C["id"],
	)
	expectedSQL(
		t,
		b,
		`SELECT "a"."value", "b"."value" FROM "relations" JOIN "a" ON "relations"."a_id" = "a"."id" JOIN "b" ON "relations"."b_id" = "b"."id"`,
		0,
	)
}
