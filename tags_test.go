package aspect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	assert := assert.New(t)

	// TODO perfect for table tests
	var tag string
	var opt options

	tag, opt = parseTag("")
	assert.Equal("", tag)
	assert.Equal(options{}, opt)
	assert.False(opt.Has("omitempty"))

	tag, opt = parseTag("id")
	assert.Equal("id", tag)
	assert.Equal(options{}, opt)
	assert.False(opt.Has("omitempty"))

	tag, opt = parseTag("id,omitempty")
	assert.Equal("id", tag)
	assert.Equal(options{"omitempty"}, opt)
	assert.True(opt.Has("omitempty"))
}
