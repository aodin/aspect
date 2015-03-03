package sqlite3

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aodin/aspect"
)

var users = aspect.Table("users",
	aspect.Column("id", aspect.Integer{NotNull: true}),
	aspect.Column("name", aspect.String{Length: 32, NotNull: true}),
	aspect.Column("password", aspect.String{Length: 128, NotNull: true}),
	aspect.PrimaryKey("id"),
)

type user struct {
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	Password string `db:"password"`
}

type id struct {
	ID int64 `db:"id"`
}

type embedded struct {
	id
	Name string `db:"name"`
}

type wrapped struct {
	ID int64 `db:"-" json:"id,omitempty"`
	user
}

// Connect to an in-memory sqlite3 instance and execute some statements.
func TestConnect(t *testing.T) {
	assert := assert.New(t)

	conn, err := aspect.Connect("sqlite3", ":memory:")
	assert.Nil(err)
	defer conn.Close()

	// Create the users table
	_, err = conn.Execute(users.Create())
	require.Nil(t, err)

	// Insert a user
	admin := user{
		ID:       1,
		Name:     "admin",
		Password: "secret",
	}
	_, err = conn.Execute(users.Insert().Values(admin))
	require.Nil(t, err, "Inserting a single user should not error")

	// Insert multiple users
	clients := []user{
		user{ID: 2, Name: "client1"},
		user{ID: 3, Name: "client2"},
	}
	_, err = conn.Execute(users.Insert().Values(clients))
	require.Nil(t, err, "Inserting multiple users should not error")

	// Select a user - Query must be given a pointer
	var u user
	require.Nil(t, conn.QueryOne(users.Select(), &u))
	assert.Equal(admin.ID, u.ID)
	assert.Equal(admin.Name, u.Name)
	assert.Equal(admin.Password, u.Password)

	// Select an incomplete, embedded struct
	var embed embedded
	require.Nil(t, conn.QueryOne(users.Select(), &embed))
	assert.Equal(admin.ID, embed.id.ID)
	assert.Equal(admin.Name, embed.Name)

	// Select a wrapped struct
	var wrap wrapped
	require.Nil(t, conn.QueryOne(users.Select(), &wrap))
	assert.Equal(admin.ID, wrap.user.ID)
	assert.Equal(admin.Name, wrap.user.Name)
	assert.Equal(admin.Password, wrap.user.Password)
	assert.Equal(0, wrap.ID)

	// Select multiple users
	var us []user
	require.Nil(t, conn.QueryAll(users.Select().OrderBy(users.C["id"]), &us))
	require.Equal(t, 3, len(us))
	assert.Equal(admin.ID, us[0].ID)
	assert.Equal(admin.Name, us[0].Name)
	assert.Equal(admin.Password, us[0].Password)

	// Select multiple users with embedding
	var embeds []embedded
	require.Nil(t,
		conn.QueryAll(users.Select().OrderBy(users.C["id"]), &embeds),
	)
	require.Equal(t, 3, len(us))
	assert.Equal(admin.ID, embeds[0].id.ID)
	assert.Equal(admin.Name, embeds[0].Name)

	// Select multiple embedded users into a slice that is pre-populated
	embeds = []embedded{
		{Name: "a"},
		{Name: "b"},
		{Name: "c"},
	}
	require.Nil(t, conn.QueryAll(
		aspect.Select(users.C["id"]).OrderBy(users.C["id"]), &embeds,
	))
	require.Equal(t, 3, len(us))
	assert.Equal(1, embeds[0].id.ID)
	assert.Equal("a", embeds[0].Name)

	// Drop the table
	_, err = conn.Execute(users.Drop())
	assert.Nil(err, "Dropping the users table should not fail")
}
