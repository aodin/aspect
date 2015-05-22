package sqlite3

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aodin/aspect"
)

// The sql dialect must implement the dialect interface
var _ aspect.Dialect = &Sqlite3{}

var things = aspect.Table("things",
	aspect.Column("name", aspect.String{Length: 32, NotNull: true}),
	aspect.Column(
		"created_at", Datetime{NotNull: true, Default: CurrentTimestamp},
	),
)

type thing struct {
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at,omitempty"`
}

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
	require.Nil(t, err, "Failed to connect to in-memory sqlite3 instance")
	defer conn.Close()

	// Create the users table
	_, err = conn.Execute(users.Create())
	require.Nil(t, err, "Failed to create users table")

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
	assert.EqualValues(0, wrap.ID)

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
	assert.EqualValues(1, embeds[0].id.ID)
	assert.Equal("a", embeds[0].Name)

	// Drop the table
	_, err = conn.Execute(users.Drop())
	assert.Nil(err, "Dropping the users table should not fail")
}

// tagless struct's exported fields must match the number and order of the
// the selected columns
type tagless struct {
	ID      int64
	Name    string
	manager interface{}
}

type fullTagless struct {
	ID       int64
	Name     string
	Password string
}

type fullname struct {
	Name string `db:"name"`
}

type embedUser struct {
	id
	fullname
	Password string `db:"password"`
}

func TestColumnTypes(t *testing.T) {
	conn, err := aspect.Connect("sqlite3", ":memory:")
	require.Nil(t, err, "Failed to connect to in-memory sqlite3 instance")
	defer conn.Close()

	// Create the things table which has a sqlite3.Datetime column type
	_, err = conn.Execute(things.Create())
	require.Nil(t, err, "Failed to create things table")

	// Sqlite3 only has per second resolution
	before := time.Now().Add(-time.Second)
	newThingy := thing{Name: "A Thing"}
	conn.MustExecute(things.Insert().Values(newThingy))
	after := time.Now().Add(time.Second)

	var thingy thing
	conn.MustQueryOne(things.Select(), &thingy)
	assert.Equal(t, "A Thing", thingy.Name)
	assert.True(t, thingy.CreatedAt.After(before))
	assert.True(t, thingy.CreatedAt.Before(after))
}

func TestResultTypes(t *testing.T) {
	conn, err := aspect.Connect("sqlite3", ":memory:")
	require.Nil(t, err, "Failed to connect to in-memory sqlite3 instance")
	defer conn.Close()

	// Create the users table
	_, err = conn.Execute(users.Create())
	require.Nil(t, err, "Failed to create users table")

	admin := embedUser{
		id:       id{ID: 1},
		fullname: fullname{Name: "admin"},
		Password: "secret",
	}
	_, err = conn.Execute(users.Insert().Values(admin))
	require.Nil(t, err, "Inserting a single user should not error")

	// Tagless destination
	var untagged tagless
	conn.MustQueryOne(
		aspect.Select(users.C["id"], users.C["name"]).Limit(1),
		&untagged,
	)
	assert.EqualValues(t, 1, untagged.ID)
	assert.Equal(t, "admin", untagged.Name)

	var untaggeds []tagless
	conn.MustQueryAll(
		aspect.Select(users.C["id"], users.C["name"]),
		&untaggeds,
	)
	require.Equal(t, 1, len(untaggeds))
	assert.EqualValues(t, 1, untaggeds[0].ID)
	assert.Equal(t, "admin", untaggeds[0].Name)

	// Tagless insert - number of columns must match numebr of exported fields
	noTag := fullTagless{
		ID:   2,
		Name: "tagless",
		// Password is a blank string
	}
	_, err = conn.Execute(users.Insert().Values(noTag))
	require.Nil(t, err, "Inserting a single tagless user should not error")

	// Embedded destination
	var embed embedUser
	conn.MustQueryOne(
		users.Select().Where(users.C["id"].Equals(1)).Limit(1),
		&embed,
	)
	assert.EqualValues(t, 1, embed.id.ID)
	assert.Equal(t, "admin", embed.fullname.Name)
	assert.Equal(t, "secret", embed.Password)
}
