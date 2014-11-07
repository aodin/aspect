package postgres

import (
	"testing"

	"github.com/aodin/aspect"
)

func TestSequence(t *testing.T) {
	expect := aspect.NewTester(t, &PostGres{})

	// Alter statements without a clause will error
	expect.Error(AlterSequence(Sequence("companies_id_seq")))
	expect.SQL(
		`ALTER SEQUENCE "companies_id_seq" RENAME TO 'whatever'`,
		AlterSequence(Sequence("companies_id_seq")).RenameTo("whatever"),
	)
	expect.SQL(
		`ALTER SEQUENCE "companies_id_seq" RESTART WITH 1`,
		AlterSequence(Sequence("companies_id_seq")).RestartWith(1),
	)
}
