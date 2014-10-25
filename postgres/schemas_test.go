package postgres

import (
	"github.com/aodin/aspect"
)

type user struct {
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	Password string `db:"password"`
	IsActive bool   `db:"is_active"`
}

var users = aspect.Table("users",
	aspect.Column("id", Serial{NotNull: true}),
	aspect.Column("name", aspect.String{Length: 32, NotNull: true}),
	aspect.Column("password", aspect.String{Length: 128}),
	aspect.Column("is_active", aspect.Boolean{Default: "TRUE"}),
	aspect.PrimaryKey("id"),
)
