package aspect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var target = Table("target",
	Column("ID", Integer{}),
	Column("A", String{}),
)

type valid struct {
	ID int64
	A  string
}

type tags struct {
	IDX int64  `db:"ID"`
	A   string `db:"a"`
}

type malformeds struct {
	IDX int64  `db:"ID`
	AX  string `db:"A`
}

type extra struct {
	ID   int64
	A, B string
}

type private struct {
	ID   int64
	A, a string
}

type missing struct {
	ID int64
	B  string
}

func TestFieldMap(t *testing.T) {
	assert := assert.New(t)

	// Target columns
	columns := target.Columns()

	var v valid
	fields, err := fieldMap(columns, v)
	// fields should map column name as key to struct field as value
	assert.Nil(err)
	assert.Equal(2, len(fields))
	assert.Equal(fields["ID"], "ID")
	assert.Equal(fields["A"], "A")

	// Should work for pointers
	fields, err = fieldMap(columns, &v)
	// fields should map column name as key to struct field as value
	assert.Nil(err)
	assert.Equal(2, len(fields))
	assert.Equal(fields["ID"], "ID")
	assert.Equal(fields["A"], "A")

	// Or tags
	var tg tags
	fields, err = fieldMap(columns, tg)
	assert.Nil(err)
	assert.Equal(2, len(fields))
	assert.Equal(fields["ID"], "IDX")
	assert.Equal(fields["A"], "A")

	// But fail for non-struct types
	var slice []int64
	_, err = fieldMap(columns, slice)
	assert.NotNil(err)

	// Ignore private fields
	var p private
	fields, err = fieldMap(columns, p)
	assert.Nil(err)
	assert.Equal(2, len(fields))
	assert.Equal(fields["ID"], "ID")
	assert.Equal(fields["A"], "A")
}
