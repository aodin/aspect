package aspect

import "testing"

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

func TestJoinOnStmt(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	// JoinOn
	expect.SQL(
		`SELECT "a"."id", "a"."value" FROM "a" JOIN "relations" ON "a"."id" = "relations"."a_id" AND "a"."id" = $1`,
		Select(tableA).JoinOn(
			relations,
			tableA.C["id"].Equals(relations.C["a_id"]),
			tableA.C["id"].Equals(2),
		),
		2,
	)

	// LeftOuterJoinOn
	expect.SQL(
		`SELECT "a"."id", "a"."value" FROM "a" LEFT OUTER JOIN "relations" ON "a"."id" = "relations"."a_id" AND "a"."id" = $1`,
		Select(tableA).LeftOuterJoinOn(
			relations,
			tableA.C["id"].Equals(relations.C["a_id"]),
			tableA.C["id"].Equals(2),
		),
		2,
	)
}
