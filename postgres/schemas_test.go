package postgres

import (
	"time"

	"github.com/aodin/aspect"
)

type user struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Password  string    `db:"password"`
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at,omitempty"`
}

var users = aspect.Table("users",
	aspect.Column("id", Serial{NotNull: true}),
	aspect.Column("name", aspect.String{Length: 32, NotNull: true}),
	aspect.Column("password", aspect.String{Length: 128}),
	aspect.Column("is_active", aspect.Boolean{Default: aspect.True}),
	aspect.Column("created_at", aspect.Timestamp{Default: Now}),
	aspect.PrimaryKey("id"),
)

type hasUUID struct {
	UUID string `db:"uuid,omitempty"`
	Name string `db:"name"`
}

var hasUUIDs = aspect.Table("has_uuids",
	aspect.Column("uuid", UUID{NotNull: true, Default: GenerateV4}),
	aspect.Column("name", aspect.String{Length: 32, NotNull: true}),
	aspect.PrimaryKey("uuid"),
)
