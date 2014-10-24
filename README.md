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

### INSERT

```go
insertUsers := Users.Insert(
    User{1, "admin", "secret"}, 
    User{2, "client", "1234"},
)
db.Execute(insertUsers)
```
```sql
INSERT INTO "users" ("id", "name", "password") VALUES ($1, $2, $3), ($4, $5, $6)
```

Structs can be partially inserted by specifying the columns. For this to work, the column names must either match the field name or the field tag `db`:

```go
type User struct {
    ID       int64  `db:"id"`
    Name     string `db:"name"`
    Password string `db:"password"`
}

admin := User{Name: "admin", Password: "secret"}
client := User{Name: "client", Password: "1234"}

insertStmt := Insert(Users.C["name"], Users.C["password"]).Values(admin, client)
db.Execute(insertStmt)
```

```sql
INSERT INTO "users" ("name", "password") VALUES ($1, $2), ($3, $4)
```


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
