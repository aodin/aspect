package aspect

import (
	"reflect"
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

type ignore struct {
	ID int64  `db:"-"`
	A  string `db:"A"`
}

type omit struct {
	ID int64  `db:"ID,omitempty"`
	A  string `db:"A"`
}

type extra struct {
	ID   int64
	A, B string
}

type private struct {
	ID   int64
	a, A string
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
	var err error
	var fields map[string]Field

	fields, err = fieldMap(columns, v)
	// fields should map column name as key to struct field as value
	assert.Nil(err)
	assert.Equal(2, len(fields))
	assert.Equal("ID", fields["ID"].Name)
	assert.Equal("A", fields["A"].Name)

	// Should work for pointers
	fields, err = fieldMap(columns, &v)
	// fields should map column name as key to struct field as value
	assert.Nil(err)
	assert.Equal(2, len(fields))
	assert.Equal("ID", fields["ID"].Name)
	assert.Equal("A", fields["A"].Name)

	// Or tags
	var tg tags
	fields, err = fieldMap(columns, tg)
	assert.Nil(err)
	assert.Equal(2, len(fields))
	assert.Equal("IDX", fields["ID"].Name)
	assert.Equal("A", fields["A"].Name)

	// But fail for non-struct types
	var slice []int64
	_, err = fieldMap(columns, slice)
	assert.NotNil(err)

	// Ignore private fields
	var p private
	fields, err = fieldMap(columns, p)
	assert.Nil(err)
	assert.Equal(2, len(fields))
	assert.Equal("ID", fields["ID"].Name)
	assert.Equal("A", fields["A"].Name)

	// Ignore fields with "-"
	var ig ignore
	fields, err = fieldMap(columns, ig)
	assert.Nil(err)
	assert.Equal(1, len(fields))
	assert.Equal("A", fields["A"].Name)

	// Omit empty fields when omitempty option is present in db struct tag
	var om omit
	fields, err = fieldMap(columns, om)
	assert.Nil(err)
	assert.Equal(2, len(fields))
	assert.Equal("ID", fields["ID"].Name)
	assert.True(fields["ID"].OmitEmpty)
	assert.Equal("A", fields["A"].Name)
}

func TestSelectAlias(t *testing.T) {
	assert := assert.New(t)

	// Determine indexes of destination struct fields
	fields := selectAlias([]string{"ID", "A"}, reflect.TypeOf(&valid{}).Elem())
	assert.Equal([]int{0, 1}, fields)

	// Determine indexes when using tags
	fields = selectAlias([]string{"ID", "A"}, reflect.TypeOf(&tags{}).Elem())
	assert.Equal([]int{0, 1}, fields)

	fields = selectAlias([]string{"ID", "A"}, reflect.TypeOf(&omit{}).Elem())
	assert.Equal([]int{0, 1}, fields)

	fields = selectAlias([]string{"ID", "A"}, reflect.TypeOf(&ignore{}).Elem())
	assert.Equal([]int{1}, fields)
}
