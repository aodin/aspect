package aspect

import (
	"testing"
)

func TestOrder(t *testing.T) {
	// Asc is implied
	o := OrderedColumn{inner: users.C["id"]}
	expectedPostGres(t, o, `"users"."id"`, 0)

	// Desc
	expectedPostGres(t, o.Desc(), `"users"."id" DESC`, 0)

	// Desc, nulls first
	expectedPostGres(
		t,
		o.Desc().NullsFirst(),
		`"users"."id" DESC NULLS FIRST`,
		0,
	)

	// Nulls last
	expectedPostGres(t, o.Asc().NullsLast(), `"users"."id" NULLS LAST`, 0)
}
