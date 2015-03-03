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

func TestFields(t *testing.T) {
	assert := assert.New(t)

	var e embedded
	fields, err := SelectFields(&e)
	require.Nil(t, err, "Fields should not error when given a proper struct")

	// It should find 4 fields: id, name, created_at, and updated_at
	require.Equal(t, 4, len(fields))

	// id
	assert.Equal(
		field{
			index:   []int{0, 0},
			column:  "id",
			table:   "",
			options: []string{"omitempty"},
		},
		fields[0],
	)

	// name
	assert.Equal(
		field{
			index:   []int{1},
			column:  "name",
			table:   "",
			options: []string{},
		},
		fields[1],
	)

	// created_at
	assert.Equal(
		field{
			index:   []int{2, 0},
			column:  "created_at",
			table:   "",
			options: []string{"omitempty"},
		},
		fields[2],
	)

	// updated_at
	assert.Equal(
		field{
			index:   []int{2, 1},
			column:  "updated_at",
			table:   "",
			options: []string{},
		},
		fields[3],
	)
}
