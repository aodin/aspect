package postgis

import (
	"github.com/aodin/aspect"
	"github.com/aodin/aspect/postgres"
	"testing"
)

var shapes = aspect.Table("shapes",
	aspect.Column("pt", Geometry{Geom: Point{}}),
	aspect.Column("area", Geometry{Polygon{}, 4326}),
)

func expectedPostGres(t *testing.T, stmt aspect.Compiles, expected string, p int) {
	params := aspect.Params()
	compiled, err := stmt.Compile(&postgres.PostGres{}, params)
	if err != nil {
		t.Error(err)
	}
	if compiled != expected {
		t.Errorf("Unexpected SQL: %s != %s", compiled, expected)
	}
	if params.Len() != p {
		t.Errorf(
			"Unexpected number of parameters for %s: %d != %d",
			expected,
			params.Len(),
			p,
		)
	}
}

func TestLatLong(t *testing.T) {
	p := LatLong{39.739167, -104.984722}
	expectedPostGres(
		t,
		p,
		`ST_SetSRID(ST_Point(-104.984722, 39.739167), 4326)::geography`,
		0,
	)
}

func TestWithin(t *testing.T) {
	expectedPostGres(
		t,
		Within(shapes.C["area"], Point{-104.984722, 39.739167}),
		`ST_Within(ST_Point(-104.984722, 39.739167), "shapes"."area")`,
		0,
	)
}

func TestGeoJSON(t *testing.T) {
	c := AsGeoJSON(shapes.C["area"])
	expectedPostGres(t, c, `ST_AsGeoJSON("shapes"."area")`, 0)
}
