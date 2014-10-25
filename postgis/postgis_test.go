package postgis

import (
	"testing"

	"github.com/aodin/aspect"
	"github.com/aodin/aspect/postgres"
)

// TODO Schemas should live in one file

var shapes = aspect.Table("shapes",
	aspect.Column("pt", Geometry{Geom: Point{}}),
	aspect.Column("area", Geometry{Polygon{}, 4326}),
)

func TestLatLong(t *testing.T) {
	expect := aspect.NewTester(t, &postgres.PostGres{})

	// TODO parameterization?
	expect.SQL(
		`ST_SetSRID(ST_Point(-104.984722, 39.739167), 4326)::geography`,
		LatLong{39.739167, -104.984722},
	)
}

func TestWithin(t *testing.T) {
	expect := aspect.NewTester(t, &postgres.PostGres{})
	expect.SQL(
		`ST_Within(ST_Point(-104.984722, 39.739167), "shapes"."area")`,
		Within(shapes.C["area"], Point{-104.984722, 39.739167}),
	)
}

func TestGeoJSON(t *testing.T) {
	expect := aspect.NewTester(t, &postgres.PostGres{})
	expect.SQL(`ST_AsGeoJSON("shapes"."area")`, AsGeoJSON(shapes.C["area"]))
}
