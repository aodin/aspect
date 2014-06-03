package aspect

import (
	"testing"
)

// TODO make admin and client globals

func testFieldIndex(t *testing.T) {
	admin := user{Id: 1, Name: "admin"}
	index, err := fieldIndex(admin, "id")
	if err != nil {
		t.Fatalf("Unexpected error during fieldIndex(): %s", err)
	}
	if index != 0 {
		t.Errorf(`Unexpected index of the user "id" column: %d`, index)
	}

	getFieldByIndex(admin, 0)
}

func TestDelete(t *testing.T) {
	// Test a complete delete
	expectedSQL(
		t,
		users.Delete(),
		`DELETE FROM "users"`,
		0,
	)

	// Test a delete with a WHERE
	expectedSQL(
		t,
		users.Delete().Where(users.C["id"].Equals(1)),
		`DELETE FROM "users" WHERE "users"."id" = $1`,
		1,
	)

	// Delete by a schema's declared primary key
	admin := user{Id: 1, Name: "admin", Password: "secret"}
	client := user{Id: 2, Name: "client", Password: "secret"}
	expectedSQL(
		t,
		users.Delete(admin),
		`DELETE FROM "users" WHERE "users"."id" = $1`,
		1,
	)

	expectedSQL(
		t,
		users.Delete(admin, client),
		`DELETE FROM "users" WHERE "users"."id" IN ($1, $2)`,
		2,
	)

	// edgeA := edge{A:1, B:2}
	// // edgeB := edge{A:2, B:3}

	// // Test delete with a composite primary key
	// expectedSQL(
	// 	t,
	// 	edges.Delete(edgeA),
	// 	`DELETE FROM "edges" WHERE "edges"."a" = $1 AND "edges"."b" = $2`,
	// 	2,
	// )
}
