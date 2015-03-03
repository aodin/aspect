package aspect

import "testing"

func TestText(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})
	expect.Create("TEXT", Text{})
}
