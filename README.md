Aspect
======

A relational database toolkit in Go.

### SQL

Executing SQL returns a results structure with methods for unpacking returned rows into structs and other data types.

```go
import (
    "fmt"
    "github.com/aodin/aspect"
    _ "github.com/lib/pq"
    "log"
)

type User struct {
    Id       int64  `db:"id"`
    Name     string `db:"name"`
    Password string `db:"password"`
}

func (u User) String() string {
    return fmt.Sprintf("%d: %s", u.Id, u.Name)
}

func main() {
    conn, err := aspect.Connect(
        "postgres",
        "host=localhost port=5432 dbname=db user=postgres password=pass",
    )
    if err != nil {
        log.Fatal("Could not connect to database")
    }
    defer conn.Close()

    // Multiple result structs
    var users []User
    if err := conn.MustExecuteSQL(`SELECT id, name, password FROM users`).All(&users); err != nil {
        log.Fatal(err)
    }
    fmt.Println(users)
    // [1: admin 2: client 3: daemon]

    // Single result struct
    var user User
    if err = conn.MustExecuteSQL(`SELECT id, name, password FROM users WHERE id = $1`, 1).One(&user); err != nil {
        log.Fatal(err)
    }
    fmt.Println(user)
    // 1: admin

    // Other result types that match the returned values
    var ids []int64
    if err = conn.MustExecuteSQL(`SELECT id FROM users`).All(&ids); err != nil {
        log.Fatal(err)
    }
    fmt.Println(ids)
    // [1 2 3]
}

```

### Schema

More advanced operations require schemas. A simple users table:

```go
import (
    * "github.com/aodin/aspect"
)

var Users = Table("users",
    Column("id", Integer{}),
    Column("name", String{"Length": 32}),
    Column("password", String{}),
    PrimaryKey("id"),
)
```

### INSERT

```go
insertUsers := Users.Insert(
    User{1, "admin", "secret"}, 
    User{2, "client", "1234"},
    User{3, "daemon", ""},
)
conn.MustExecute(insertUsers)
```
```sql
INSERT INTO "users" ("id", "name", "password") VALUES ($1, $2, $3), ($4, $5, $6), ($7, $8, $9)
```

Structs can be partially inserted by specifying the columns. For this to work, the column names must either match the field name or the field tag `db`:

```go
type User struct {
    Id       int64  `db:"id"`
    Name     string `db:"name"`
    Password string `db:"password"`
}

admin := user{Name: "admin", Password: "secret"}
client := user{Name: "client", Password: "1234"}

insertColumns := Insert(Users.C["name"], Users.C["password"])
insertColumns = insertColumns.Values(admin, client)
conn.MustExecute(insertColumns)
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
result, err := conn.Execute(Users.Select())
if err != nil {
    result.All(&users)
}
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
conn.MustExecute(s).All(&ids)
fmt.Println(ids)
// [3, 2, 1]
```

### Delete

```go
Users.Delete()
```

```sql
DELETE FROM "users"
```

> Death and Light are everywhere, always, and they begin, end, strive,
> attend, into and upon the Dream of the Nameless that is the world,
> burning words within Samsara, perhaps to create a thing of beauty.
>
> _Lord of Light_ by Roger Zelazny
