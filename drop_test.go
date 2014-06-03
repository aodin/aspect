package aspect

import (
	"testing"
)

func TestDropStmt(t *testing.T) {
	expectedSQL(t, users.Drop(), `DROP TABLE "users"`, 0)
}
