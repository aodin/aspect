package aspect

import (
	"fmt"
)

type Deletable interface {
	Deletable() *TableStruct
}

type DeleteStatement struct {
	Target      Deletable
}

func (stmt *DeleteStatement) String() string {
	return stmt.Compile()
}

func (stmt *DeleteStatement) Compile() string {
	return fmt.Sprintf(`DELETE FROM "%s"`, stmt.Target.Deletable().Name)
}

func (stmt *DeleteStatement) Execute() (string, error) {
	// TODO Return any delayed errors
	// TODO Check for a cached string
	return stmt.Compile(), nil
}
