package aspect

import (
	"testing"
)

func TestColumnElem(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	// Calling inner should return the inner ColumnClause
	col := users.C["id"]
	copy := col.Inner()
	expect.SQL(`"users"."id"`, copy)

	// Change the inner
	resetCol := col.SetInner(copy)
	expect.SQL(`"users"."id"`, resetCol)

	// Test the old compilation behavior
	chad := ColumnElem{
		table: users,
		name:  "chad",
	}
	expect.SQL(`"users"."chad"`, chad)
}

func TestColumnAlias(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})
	expect.SQL(
		`"users"."name" AS "user_name"`,
		users.C["name"].As("user_name"),
	)
}

func TestColumnOrdering(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})
	expect.SQL(`"users"."id"`, users.C["id"].Orderable())
	expect.SQL(`"users"."id"`, users.C["id"].Asc())
	expect.SQL(`"users"."id" DESC`, users.C["id"].Desc())
	expect.SQL(`"users"."id" NULLS FIRST`, users.C["id"].NullsFirst())
	expect.SQL(`"users"."id" NULLS LAST`, users.C["id"].NullsLast())
}

func TestColumnConditionals(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})
	expect.SQL(`"users"."id" = $1`, users.C["id"].Equals(1), 1)
	expect.SQL(`"users"."id" != $1`, users.C["id"].DoesNotEqual(1), 1)
	expect.SQL(`"users"."id" < $1`, users.C["id"].LessThan(1), 1)
	expect.SQL(`"users"."id" > $1`, users.C["id"].GreaterThan(1), 1)
	expect.SQL(`"users"."id" <= $1`, users.C["id"].LTE(1), 1)
	expect.SQL(`"users"."id" >= $1`, users.C["id"].GTE(1), 1)
	expect.SQL(`"users"."name" LIKE $1`, users.C["name"].Like(`_b`), "_b")
	expect.SQL(
		`"users"."name" NOT LIKE $1`,
		users.C["name"].NotLike(`_b`),
		"_b",
	)
	expect.SQL(`"users"."name" ILIKE $1`, users.C["name"].ILike(`_b`), "_b")
	expect.SQL(
		`"users"."name" SIMILAR TO $1`,
		users.C["name"].SimilarTo(`_b`),
		"_b",
	)
	expect.SQL(
		`"users"."name" NOT SIMILAR TO $1`,
		users.C["name"].NotSimilarTo(`_b`),
		"_b",
	)
	expect.SQL(`"users"."id" IS NULL`, users.C["id"].IsNull())
	expect.SQL(`"users"."id" IS NOT NULL`, users.C["id"].IsNotNull())
}

func TestColumnElem_Modify(t *testing.T) {
	table := &TableElem{Name: "users"}

	// Attempt to modify a table with a nameless columns
	nameless := ColumnElem{}
	if err := nameless.Modify(table); err == nil {
		t.Fatalf("no error when modifying a table with a nameless column")
	}

	// Attempt to modify a nil table
	named := ColumnElem{name: "id"}
	if err := named.Modify(nil); err == nil {
		t.Fatalf("no error when modifying a nil table")
	}

	// Add a column to the same table twice
	twice := ColumnElem{name: "id", table: table}
	if err := twice.Modify(table); err == nil {
		t.Fatalf("no error when adding a column which already has a table")
	}
}
