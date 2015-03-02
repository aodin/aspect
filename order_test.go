package aspect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrder(t *testing.T) {
	assert := assert.New(t)
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
	assert.Equal(o.inner.Name(), o.Orderable().inner.Name())
}
