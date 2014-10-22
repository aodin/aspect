package aspect

import (
	"testing"
)

func TestOrder(t *testing.T) {
	// Asc is implied
	o := OrderedColumn{inner: users.C["id"]}
	expectedSQL(t, o, `"users"."id"`, 0)

	// Desc
	expectedSQL(t, o.Desc(), `"users"."id" DESC`, 0)

	// Desc, nulls first
	expectedSQL(
		t,
		o.Desc().NullsFirst(),
		`"users"."id" DESC NULLS FIRST`,
		0,
	)

	// Asc, Nulls last
	expectedSQL(t, o.Asc().NullsLast(), `"users"."id" NULLS LAST`, 0)

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
