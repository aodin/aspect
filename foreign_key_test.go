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

var messages = Table("messages",
	Column("id", Integer{NotNull: true, PrimaryKey: true}),
	SelfForeignKey("parent_id", "id", Integer{}),
)

func TestSelfForeignKey(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	expect.SQL(
		`CREATE TABLE "messages" (
  "id" INTEGER PRIMARY KEY NOT NULL,
  "parent_id" INTEGER REFERENCES messages("id")
);`,
		messages.Create(),
	)
}

func TestForeignKeyElement(t *testing.T) {
	assert := assert.New(t)

	// Test that the fk columns were added to the C mapping
	_, ok := children.C["parent_id"]
	assert.Equal(true, ok)

	// There should also be a foreign key element attached to the table
	fks := children.ForeignKeys()
	assert.Equal(1, len(fks))

	fk := fks[0]
	assert.Equal(children, fk.Table())
	assert.Equal(parents, fk.ReferencesTable())
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
