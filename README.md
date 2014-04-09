Aspect
======

A relational database toolkit in Go.

### Schema

A simple table schema:

```go
import (
    "github.com/aodin/aspect"
)

var Users = aspect.Table("users",
    aspect.Column("id", aspect.Integer{}),
    aspect.Column("name", aspect.String{"Length": 32}),
    aspect.Column("password", aspect.String{}),
    aspect.PrimaryKey("id"),
)
```

### Select

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
type User struct {
    Id       int64
    Name     string
    Password string
}

func (u User) String() string {
    return fmt.Sprintf("%d: %s", u.Id, u.Name)
}

var users []User
result, err := conn.Execute(Users.Select())
if err != nil {
    result.All(&users)
}
for _, user := range users {
    fmt.Println(user)    
}
```

Simple queries can be returned into more concise types:

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
