package aspect

import (
	"testing"
)

func TestOrder(t *testing.T) {
	// Since OrderedColumn is a pointer, settings will persist
	// Asc is implied
	o := &OrderedColumn{inner: users.C["id"]}
	expectedSQL(t, o, `"users"."id"`)

	// Desc
	expectedSQL(t, o.Desc(), `"users"."id" DESC`)

	// Desc, nulls first
	expectedSQL(t, o.NullsFirst(), `"users"."id" DESC NULLS FIRST`)

	// nulls last
	expectedSQL(t, o.Asc().NullsLast(), `"users"."id" NULLS LAST`)
}
