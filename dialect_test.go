package aspect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistry(t *testing.T) {
	// Add dialects
	RegisterDialect("default", &defaultDialect{})

	assert.Panics(t,
		func() { RegisterDialect("nil", nil) },
		"registry failed to panic when given a nil dialect",
	)

	assert.Panics(t,
		func() { RegisterDialect("default", &defaultDialect{}) },
		"registry failed to panic when given a duplicate dialect",
	)

	// Get dialects
	var err error
	_, err = GetDialect("default")
	assert.Nil(t, err, "Getting a dialect that exists should not error")

	_, err = GetDialect("dne")
	assert.NotNil(t, err, "Getting a dialect that does not exist should error")
}
