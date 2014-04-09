Aspect
======

A relational database toolkit in Go.

### Schema

A simple table schema:

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

### SELECT

Each of the following statements will produce the same SQL:

```go
users.Select()
sql.Select(users)
sql.Select(users.C["id"], users.C["name"], users.C["password"])
```

```sql
SELECT "users"."id", "users"."name", "users"."password" FROM "users"
```

Results can be returned directly into structs:

```go
type User struct {
    Id       int64
    Name     string
    Password string
}

func (u User) String() string {
    return fmt.Sprintf("%d: %s", u.Id, u.Name)
}

var users []User
result, err := conn.Execute(users.Select())
if err != nil {
    result.All(&users)
}
for _, user := range users {
    fmt.Println(user)    
}
```

And simplier queries into more concise return types:

```go
s := aspect.Select(Users.C["id"]).OrderBy(Users.C["id"])
fmt.Println(s)
```

```sql
SELECT "users"."id" FROM "users" ORDER BY "users"."id"
```

```go
var ids []int64
conn.MustExecute(s).All(&ids)
fmt.Println(ids)
```

> Death and Light are everywhere, always, and they begin, end, strive,
> attend, into and upon the Dream of the Nameless that is the world,
> burning words within Samsara, perhaps to create a thing of beauty.
>
> _Lord of Light_ by Roger Zelazny
