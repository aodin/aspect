package aspect

import (
	"testing"
)

func TestDropStmt(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})
	expect.SQL(`DROP TABLE "users"`, users.Drop())
}
