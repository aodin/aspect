package aspect

import (
	"testing"
)

func TestOrder(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	// Asc is implied
	o := OrderedColumn{inner: users.C["id"]}
	expect.SQL(`"users"."id"`, o)

	// Desc
	expect.SQL(`"users"."id" DESC`, o.Desc())

	// Desc, nulls first
	expect.SQL(
		`"users"."id" DESC NULLS FIRST`,
		o.Desc().NullsFirst(),
	)

	// Asc, Nulls last
	expect.SQL(`"users"."id" NULLS LAST`, o.Asc().NullsLast())

	// Calling Orderable on an OrderableColumn should return a copy of itself
	copy := o.Orderable()
	if copy.inner.Name() != o.inner.Name() {
		t.Fatal(
			"unexpected inner name during orderable copy: %s != %s",
			copy.inner.Name(),
			o.inner.Name(),
		)
	}
}
