package postgres

import (
	"testing"

	"github.com/aodin/aspect"
)

func TestColumn(t *testing.T) {
	expect := aspect.NewTester(t, &PostGres{})
	expect.SQL(
		`"times"."when" @> $1`,
		C(times.C["when"]).Contains("[0, 1)"),
		"[0, 1)",
	)

	// JSON selector
	// http://www.postgresql.org/docs/9.3/static/functions-json.html
	expect.SQL(
		`"members"."info" ->> 'name' AS "Name"`,
		C(members.C["info"]).GetJSONText("name").As("Name"),
	)
}
