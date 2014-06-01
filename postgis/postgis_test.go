package postgis

import (
	"github.com/aodin/aspect"
	"testing"
)

var shapes = aspect.Table("shapes",
	aspect.Column("geom", GeometryPolygon{4326}),
)

func expectedPostGres(t *testing.T, stmt aspect.Compiles, expected string, p int) {
	params := aspect.Params()
	compiled, err := stmt.Compile(&aspect.PostGres{}, params)
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

func TestPoint(t *testing.T) {
	p := Point{39.739167, -104.984722}
	expectedPostGres(
		t,
		p,
		`ST_GeometryFromText('POINT(-104.984722 39.739167)', 4326)`,
		0,
	)
}

func TestWithin(t *testing.T) {
	expectedPostGres(
		t,
		Within(shapes.C["geom"], Point{39.739167, -104.984722}),
		`ST_Within(ST_GeometryFromText('POINT(-104.984722 39.739167)', 4326), "shapes"."geom")`,
		0,
	)
}

func TestGeoJSON(t *testing.T) {
	c := GeoJSON(shapes.C["geom"])
	expectedPostGres(t, c, `ST_AsGeoJSON("shapes"."geom")`, 0)
}
