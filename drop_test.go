package aspect

import (
	"testing"
)

func TestDropStmt(t *testing.T) {
	expectedPostGres(t, users.Drop(), `DROP TABLE "users"`, 0)
}
