package aspect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var clients = Table("clients",
	Column("id", Integer{NotNull: true}),
	Column("name", String{Length: 32, Unique: true, NotNull: true}),
	Column("password", String{Length: 128}),
	PrimaryKey("id"),
)

type client struct {
	ID       int64  `db:"id,omitmpty"`
	Name     string `db:"name"`
	Password string `db:"password"`
}

func TestValidateInsert(t *testing.T) {
	assert := assert.New(t)

	// A zero-init client is valid
	c := client{}
	err := ValidateInsert(clients, c)
	assert.Nil(err)
}
