package aspect

import (
	"testing"
)

func TestUpdate(t *testing.T) {
	s1 := Update(users, Values{"name": "client"})
	expectedSQL(
		t,
		s1,
		`UPDATE "users" SET "name" = $1`,
		1,
	)

	values := Values{
		"name":     "admin",
		"password": "blank",
	}

	s2 := Update(users, values).Where(users.C["id"].Equals(1))
	expectedSQL(
		t,
		s2,
		`UPDATE "users" SET "name" = $1 AND "password" = $2 WHERE "users"."id" = $3`,
		3,
	)

	// The statement should have an error if a values key does not have an
	// associated column
	s3 := Update(users, Values{})
	_, err := s3.Compile(&defaultDialect{}, Params())
	if err == nil {
		t.Fatalf("No error returned from column-less UPDATE")
	}
}
