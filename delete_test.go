package aspect

import (
	"testing"
)

func TestFieldIndex(t *testing.T) {
	admin := user{ID: 1, Name: "admin"}
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
	expect := NewTester(t, &defaultDialect{})

	// Test a complete delete
	expect.SQL(
		`DELETE FROM "users"`,
		users.Delete(),
	)

	// Test a delete with a WHERE
	expect.SQL(
		`DELETE FROM "users" WHERE "users"."id" = $1`,
		users.Delete().Where(users.C["id"].Equals(1)),
		1,
	)

	// Delete by a schema's declared primary key
	admin := user{ID: 1, Name: "admin", Password: "secret"}
	client := user{ID: 2, Name: "client", Password: "secret"}
	expect.SQL(
		`DELETE FROM "users" WHERE "users"."id" = $1`,
		users.Delete(admin),
		1,
	)

	expect.SQL(
		`DELETE FROM "users" WHERE "users"."id" IN ($1, $2)`,
		users.Delete(admin, client),
		1,
		2,
	)

}
