package aspect

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type embeddedID struct {
	ID int64 `db:"id,omitempty"`
}

type embedded struct {
	embeddedID
	Name      string `db:"name"`
	Timestamp struct {
		CreatedAt time.Time  `db:"created_at,omitempty"`
		UpdatedAt *time.Time `db:"updated_at"`
		is_active bool
	}
	manager *struct{}
}

type nested struct {
	Another string `db:"another"`
	embedded
}

type moreNesting struct {
	nested
	OneMore string `db:"one_more"`
}

func TestFields(t *testing.T) {
	assert := assert.New(t)

	var e embedded
	fields := SelectFields(&e)

	// It should find 4 fields: id, name, created_at, and updated_at
	require.Equal(t, 4, len(fields))

	// id
	assert.Equal(
		field{
			index:   []int{0, 0},
			column:  "id",
			options: []string{"omitempty"},
		},
		fields[0],
	)

	// name
	assert.Equal(
		field{
			index:   []int{1},
			column:  "name",
			options: []string{},
		},
		fields[1],
	)

	// created_at
	assert.Equal(
		field{
			index:   []int{2, 0},
			column:  "created_at",
			options: []string{"omitempty"},
		},
		fields[2],
	)

	// updated_at
	assert.Equal(
		field{
			index:   []int{2, 1},
			column:  "updated_at",
			options: []string{},
		},
		fields[3],
	)

	// Test a deeply nested struct
	var n moreNesting
	fields = SelectFields(&n)
	require.Equal(t, 6, len(fields))
	assert.Equal(
		field{
			index:   []int{0, 1, 2, 1},
			column:  "updated_at",
			options: []string{},
		},
		fields[4],
	)
}

func TestAlignColumns(t *testing.T) {
	assert := assert.New(t)

	fields := []field{
		field{
			index:  []int{1},
			column: "name",
		},
		field{
			index:  []int{0},
			column: "id",
		},
	}

	columns := []string{"id", "name", "age"}
	aligned := AlignColumns(columns, fields)
	require.Equal(t, 3, len(aligned))

	assert.Equal(fields[1], aligned[0])
	assert.Equal(fields[0], aligned[1])
	assert.False(aligned[2].Exists())
}
