package aspect

import (
	"testing"
)

func TestColumnElem(t *testing.T) {
	// Calling inner should return the inner ColumnClause
	col := users.C["id"]
	copy := col.Inner()
	expectedSQL(t, copy, `"users"."id"`, 0)

	// Change the inner
	resetCol := col.SetInner(copy)
	expectedSQL(t, resetCol, `"users"."id"`, 0)

	// Test the old compilation behavior
	chad := ColumnElem{
		table: users,
		name:  "chad",
	}
	expectedSQL(t, chad, `"users"."chad"`, 0)
}

func TestColumnOrdering(t *testing.T) {
	expectedSQL(t, users.C["id"].Orderable(), `"users"."id"`, 0)
	expectedSQL(t, users.C["id"].Asc(), `"users"."id"`, 0)
	expectedSQL(t, users.C["id"].Desc(), `"users"."id" DESC`, 0)
	expectedSQL(t, users.C["id"].NullsFirst(), `"users"."id" NULLS FIRST`, 0)
	expectedSQL(t, users.C["id"].NullsLast(), `"users"."id" NULLS LAST`, 0)
}

func TestColumnConditionals(t *testing.T) {
	expectedSQL(t, users.C["id"].Equals(1), `"users"."id" = $1`, 1)
	expectedSQL(t, users.C["id"].DoesNotEqual(1), `"users"."id" != $1`, 1)
	expectedSQL(t, users.C["id"].LessThan(1), `"users"."id" < $1`, 1)
	expectedSQL(t, users.C["id"].GreaterThan(1), `"users"."id" > $1`, 1)
	expectedSQL(t, users.C["id"].LTE(1), `"users"."id" <= $1`, 1)
	expectedSQL(t, users.C["id"].GTE(1), `"users"."id" >= $1`, 1)
	expectedSQL(t, users.C["id"].IsNull(), `"users"."id" IS NULL`, 0)
	expectedSQL(t, users.C["id"].IsNotNull(), `"users"."id" IS NOT NULL`, 0)
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
