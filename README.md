Aspect [![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/aodin/aspect) [![Build Status](https://travis-ci.org/aodin/aspect.svg?branch=master)](https://travis-ci.org/aodin/aspect)
======

A relational database toolkit for Go:

* Build complete database schemas
* Create reusable and cross-dialect SQL statements
* Allow struct instances and slices to be directly populated by the database

Quickstart
----------

```go
package main

import (
    "log"

    sql "github.com/aodin/aspect"
    _ "github.com/aodin/aspect/sqlite3"
)

// Create a database schema using aspect's Table function
var Users = sql.Table("users",
    sql.Column("id", sql.Integer{NotNull: true}),
    sql.Column("name", sql.String{Length: 32, NotNull: true}),
    sql.Column("password", sql.String{Length: 128}),
    sql.PrimaryKey("id"),
)

// Structs are used to send and receive values to the database
type User struct {
    ID       int64  `db:"id"`
    Name     string `db:"name"`
    Password string `db:"password"`
}

func main() {
    // Connect to an in-memory sqlite3 instance
    conn, err := sql.Connect("sqlite3", ":memory:")
    if err != nil {
        log.Panic(err)
    }
    defer conn.Close()

    // Create the users table
    conn.MustExecute(Users.Create())

    // Insert a user - they can be inserted by value or reference
    admin := User{ID: 1, Name: "admin", Password: "secret"}
    conn.MustExecute(Users.Insert().Values(admin))

    // Select a user - query methods must be given a pointer
    var user User
    conn.MustQueryOne(Users.Select(), &user)
    log.Println(user)
}
```

Statements
----------

Don't forget to import aspect and at least one driver. I alias the aspect package to `sql`:

```go
import (
    sql "github.com/aodin/aspect"
    _ "github.com/aodin/aspect/postgres"
    _ "github.com/aodin/aspect/sqlite3"
)
```

Statements that do not return selections can be run with the `Execute` method of database connections `DB` or transactions `TX`. Both also implement the interface `Connection`.

A successful `Connect` will return a database connection pool ready for use. Its `Execute` method returns an instance of `database/sql` package's `Result` and an error if one occurred:

```go
conn, err := sql.Connect("sqlite3", ":memory:")
if err != nil {
    log.Panic(err)
}
defer conn.Close()

result, err := conn.Execute(Users.Create())
```

Results are often ignored, as in the `Quickstart` example above.

The following commands are usually used with the `Execute` method:


#### CREATE TABLE

Once a schema has been specified with `Table`, such as:

```go
var Users = sql.Table("users",
    sql.Column("id", sql.Integer{NotNull: true}),
    sql.Column("name", sql.String{Length: 32, NotNull: true}),
    sql.Column("password", sql.String{Length: 128}),
    sql.PrimaryKey("id"),
)
```

A `CREATE TABLE` statement can be created with:

```go
Users.Create()
```

And will output the following SQL with its `String()` method (a dialect neutral version) or with `conn.String()` (a dialect specific version):

```sql
CREATE TABLE "users" (
  "id" INTEGER NOT NULL,
  "name" VARCHAR(32) NOT NULL,
  "password" VARCHAR(128),
  PRIMARY KEY ("id")
);
```

#### DROP TABLE

Using the `Users` schema, a `DROP TABLE` statement can be created with:

```go
Users.Drop()
```

```sql
DROP TABLE "users"
```

#### INSERT

Insert statements can be created without specifying values. For instance, the method `Insert()` on a schema such as `Users` can be created with:

```go
Users.Insert()
```

And produces the SQL (in this example for the `sqlite3` dialect):

```sql
INSERT INTO "users" ("id", "name", "password") VALUES (?, ?, ?)
```

Values can be inserted to the database using structs or `aspect.Values` instances. If given a struct, Aspect first attempts to match field names or `db` tags. Columns without matching values will be dropped. The following struct and chained function:

```go
type user struct {
    Name     string `db:"name"`
    Password string `db:"password"`
    Extra    string
    manager  *manager
}
```

```go
Users.Insert().Values(user{Name: "Totti", Password: "GOAL"})
```

```sql
INSERT INTO "users" ("name", "password") VALUES (?, ?)
```

Fields must be exported (i.e. start with an uppercase character) to work with Aspect. Unexported fields and those that do not match column names will be ignored.

Structs without `db` tags or fields that match column names can be inserted, but only if the number of exported fields equals the number of columns being inserted.

Slices of structs can also be inserted:

```go
users := []user{
    {Name: "Howard", Password: "DENIED"},
    {Name: "Beckham", Password: "RETIRED"},
}
```

And would produce the following (note: this syntax for multiple inserts is only valid in later versions of `sqlite3`):

```sql
INSERT INTO "users" ("name", "password") VALUES (?, ?), (?, ?)
```

To manually specify which columns should be inserted, use Aspect's `Insert` function, rather than the table method:

```go
sql.Insert(Users.C["name"]).Values(users)
```

```sql
INSERT INTO "users" ("name") VALUES (?), (?)
```

Values can also be inserted using the map type `Values` or a slice of them:

```go
Users.Insert().Values(sql.Values{"name": "Ronaldo"})
```

```sql
INSERT INTO "users" ("name") VALUES (?)
```

Keys in `Values` maps must match column names or the statement will error.


#### UPDATE

Rows in a table can be updated using either of:

```go
Users.Update()
sql.Update(Users)
```

Both will produce the same output:

```sql
UPDATE "users" SET "id" = ?, "name" = ?, "password" = ?
```

Columns are set according to the values given. Values keys that do not match a column will error.

```go
Users.Update().Values(sql.Values{"name": "Ronaldo", "password": "STRIKER"})
```

```sql
UPDATE "users" SET "name" = ?, "password" = ?
```

Conditionals can be specified with the `Where()` method:

```go
Users.Update().Values(
    sql.Values{"password": "FIFA2014"},
).Where(Users.C["name"].Equals("Ronaldo"))
```

```sql
UPDATE "users" SET "password" = ? WHERE "users"."name" = ?
```


#### DELETE

If you want to delete all the rows in a table:

```go
Users.Delete()
```

```sql
DELETE FROM "users"
```

A conditional delete can be created with `Where()`:

```go
Users.Delete().Where(Users.C["name"].Equals("Ronaldo"))
```

```sql
DELETE FROM "users" WHERE "users"."name" = ?
```

If the schema has a single column primary key specified, deletes can be performed with structs:

```go
Users.Delete().Values(user{ID: 1})
```

```sql
DELETE FROM "users" WHERE "users"."id" = ?
```

Or slices of structs:

```go
Users.Delete().Values([]user{{ID: 1}, {ID: 2}})
```

```sql
DELETE FROM "users" WHERE "users"."id" IN (?, ?)
```


#### SELECT

Results can be queried in a number of ways. Each of the following statements will produce the same SQL output:

```go
Users.Select()
sql.Select(Users)
sql.Select(Users.C["id"], Users.C["name"], Users.C["password"])
```

```sql
SELECT "users"."id", "users"."name", "users"."password" FROM "users"
```

Results can be returned directly into structs:

```go
var users []User
conn.MustQueryAll(Users.Select(), &users)
log.Println(users)
// [1: admin 2: client 3: daemon]
```

It's okay to have a destination struct larger than the expected results (vice-versa!), as long as the struct has `db` tags:

```go
type username struct {
    Name    string `db:"name"`
    manager *manager
}

var usernames []username
conn.MustQueryAll(sql.Select(Users.C["name"]), &usernames)
```

Single column queries can be returned into slice types:

```go
stmt := sql.Select(Users.C["id"]).OrderBy(Users.C["id"].Desc())
```

```sql
SELECT "users"."id" FROM "users" ORDER BY "users"."id" DESC
```

```go
var ids []int64
conn.MustQueryAll(s, &ids)
log.Println(ids)
// [3, 2, 1]
```

Single results can also be queried into instantiated instances of `Values`, which is included with the package:

```go
result := sql.Values{}
conn.MustQueryOne(Users.Select(), result)
```

Or multiple results with a pointer to a slice of `Values`:

```go
var results []sql.Values
conn.QueryAll(Users.Select(), &results)
```

A warning, however: string results are added to the map as `[]byte` types. Since this is done internally by the `database/sql` standard library package, this is unlikely to change in the future.

    map[id:1 name:[84 111 116 116 105] password:[71 79 65 76]]

To get string output, simply cast the value to a `string` after asserting that it is a `[]byte`:

```go
log.Println(string(result["name"].([]byte)))
// Totti
```


Schema
------

More advanced schemas can be created with foreign keys, unique constraints, and composite primary keys:

```go
var Posts = sql.Table("posts",
    sql.ForeignKey("uid", Users.C["id"], sql.BigInt{}).OnDelete(sql.Cascade),
    sql.Column("name", sql.String{Length: 32, NotNull: true}),
    sql.Column("is_published", sql.Boolean{Default: sql.True}),
    sql.Unique("name"),
    sql.PrimaryKey("uid", "name"),
)
```

```sql
CREATE TABLE "posts" (
  "uid" BIGINT REFERENCES users("id") ON DELETE CASCADE,
  "name" VARCHAR(32) NOT NULL,
  "is_published" BOOLEAN DEFAULT TRUE,
  UNIQUE ("name"),
  PRIMARY KEY ("uid", "name")
);
```


Development
-----------

To perform tests, copy the `db.example.json` file into each driver package requiring credentials as `db.json` (currently required only by the `postgres` package) and set it for your local configuration.

All compilable statements implement a `String` method for dialect-neutral logging and the `Compile` method for building dialect-specific and parameterized output:

```go
type Compiles interface {
    String() string
    Compile(Dialect, *Parameters) (string, error)
}
```

Statements and clauses can be tested by creating a new dialect-specific tester; for example using the `postgres` package:

```go
expect := NewTester(t, &postgres.PostGres{})
```

The instance's `SQL` method will test expected output and parameterization:

```go
expect.SQL(`DELETE FROM "users"`, users.Delete())

expect.SQL(
    `INSERT INTO "users" ("name") VALUES ($1), ($2)`,
    users.Insert().Values([]sql.Values{{"name": "Totti"}, {"name": "De Rossi"}}),
    "Totti",
    "De Rossi",
)
```

And the `Error` method will test that an error occurred:

```go
expect.Error(sql.Select(users.C["does-not-exist"]))
```


Happy Hacking!

aodin, 2014-15

> Death and Light are everywhere, always, and they begin, end, strive,
> attend, into and upon the Dream of the Nameless that is the world,
> burning words within Samsara, perhaps to create a thing of beauty.
>
> _Lord of Light_ by Roger Zelazny
