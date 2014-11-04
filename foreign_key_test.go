package aspect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var parents = Table("parents",
	Column("id", Integer{}),
	Column("name", String{}),
)

var children = Table("children",
	ForeignKey("parent_id", parents.C["id"]),
	Column("name", String{}),
)

var childrenType = Table("children",
	ForeignKey("p_id", parents.C["id"], BigInt{NotNull: true}),
)

var childrenCascade = Table("children",
	ForeignKey("p_id", parents.C["id"]).OnDelete(Cascade).OnUpdate(Cascade),
)

func TestForeignKey(t *testing.T) {
	// Test that the fk columns were added to the C mapping
	_, ok := children.C["parent_id"]
	assert.Equal(t, true, ok)
}

func TestForeignKey_Create(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	var expected string
	expected = `CREATE TABLE "children" (
  "parent_id" INTEGER REFERENCES parents("id"),
  "name" VARCHAR
);`
	expect.SQL(expected, children.Create())

	// Override the type of the foreign key
	expected = `CREATE TABLE "children" (
  "p_id" BIGINT NOT NULL REFERENCES parents("id")
);`
	expect.SQL(expected, childrenType.Create())

	// Add cascade behavior
	expected = `CREATE TABLE "children" (
  "p_id" INTEGER REFERENCES parents("id") ON DELETE CASCADE ON UPDATE CASCADE
);`
	expect.SQL(expected, childrenCascade.Create())

	// Test too many overrides
	func() {
		defer func() {
			if panicked := recover(); panicked == nil {
				t.Errorf("table failed to panic when multiple overriding types were added to a foreign key")
			}
		}()
		Table("bad",
			ForeignKey("no", parents.C["id"], String{}, Integer{}),
		)
	}()

}
