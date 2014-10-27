Aspect
======

A relational database toolkit in Go that aims to:

* Build complete database schemas
* Create reusable and cross-dialect SQL statements
* Allow struct instances and slices to be directly populated by the database

### Quickstart

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
    db, err := sql.Connect("sqlite3", ":memory:")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    // Create the users table
    if _, err = db.Execute(Users.Create()); err != nil {
        panic(err)
    }

    // Insert a user - they can be inserted by value or reference
    admin := User{ID: 1, Name: "admin", Password: "secret"}
    if _, err = db.Execute(Users.Insert(admin)); err != nil {
        panic(err)
    }

    // Select a user - query methods must be given a pointer
    var user User
    if err = db.QueryOne(Users.Select(), &user); err != nil {
        panic(err)
    }
    log.Println(user)
}
```

Example Statements
------------------

Don't forget to import aspect and at least one driver you'll be using. I often alias the aspect package to `sql` as below:

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
db, err := sql.Connect("sqlite3", ":memory:")
if err != nil {
    panic(err)
}
defer db.Close()

result, err := db.Execute(Users.Create())
```

Results are often ignored, as in the `Quickstart` example above.

The following commands are usually used with the `Execute` method:


### CREATE TABLE

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

And will output the following SQL with its `String()` method (a dialect neutral version) or with `db.String()` (a dialect specific version):

```sql
CREATE TABLE "users" (
  "id" INTEGER NOT NULL,
  "name" VARCHAR(32) NOT NULL,
  "password" VARCHAR(128),
  PRIMARY KEY ("id")
);
```

### DROP TABLE

Using the `Users` schema, a `DROP TABLE` statement can be created with:

```go
Users.Drop()
```

And produces the SQL:

```sql
DROP TABLE "users"
```

### INSERT

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

Will produce:

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


### SELECT

Each of the following statements will produce the same SQL:

```go
Users.Select()
aspect.Select(Users)
aspect.Select(Users.C["id"], Users.C["name"], Users.C["password"])
```

```sql
SELECT "users"."id", "users"."name", "users"."password" FROM "users"
```

Results can be returned directly into structs:

```go
var users []User
err := db.QueryAll(Users.Select(), &users)
fmt.Println(users)
// [1: admin 2: client 3: daemon]
```

Simple queries can be returned into more concise types:

```go
s := aspect.Select(Users.C["id"]).OrderBy(Users.C["id"].Desc())
```

```sql
SELECT "users"."id" FROM "users" ORDER BY "users"."id" DESC
```

```go
var ids []int64
if err := db.QueryAll(s, &ids); err != nil {
    log.Fatal(err)
}
fmt.Println(ids)
// [3, 2, 1]
```

### DELETE

```go
Users.Delete()
```

```sql
DELETE FROM "users"
```

If the schema has a primary key specified, deletes can be performed with structs:

```go
admin = User{1, "admin", "secret"}
Users.Delete(admin)
```

```sql
DELETE FROM "users" WHERE "users"."id" = $1
```

> Death and Light are everywhere, always, and they begin, end, strive,
> attend, into and upon the Dream of the Nameless that is the world,
> burning words within Samsara, perhaps to create a thing of beauty.
>
> _Lord of Light_ by Roger Zelazny
