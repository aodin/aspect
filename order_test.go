package aspect

import (
	"testing"
)

func TestOrder(t *testing.T) {
	// Since OrderedColumn is a pointer, settings will persist
	// Asc is implied
	o := &OrderedColumn{inner: users.C["id"]}
	expectedPostGres(t, o, `"users"."id"`, 0)

	// Desc
	expectedPostGres(t, o.Desc(), `"users"."id" DESC`, 0)

	// Desc, nulls first
	expectedPostGres(t, o.NullsFirst(), `"users"."id" DESC NULLS FIRST`, 0)

	// nulls last
	expectedPostGres(t, o.Asc().NullsLast(), `"users"."id" NULLS LAST`, 0)
}
