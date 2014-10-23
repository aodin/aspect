package postgis

import (
	"fmt"
	"github.com/aodin/aspect"
	"github.com/aodin/aspect/postgres"
)

// Importing PostGIS implies you'll be using PostGres

type LatLong struct {
	Latitude, Longitude float64
}

// String returns the element's SQL using the default dialect.
func (p LatLong) String() string {
	compiled, _ := p.Compile(&postgres.PostGres{}, aspect.Params())
	return compiled
}

// A point with the implied SRID of 4326
// TODO parameterization
func (p LatLong) Compile(d aspect.Dialect, params *aspect.Parameters) (string, error) {
	return fmt.Sprintf(
		`ST_SetSRID(ST_Point(%f, %f), 4326)::geography`,
		p.Longitude,
		p.Latitude,
	), nil
}

func (p LatLong) Create(d aspect.Dialect) (string, error) {
	return "POINT", nil
}

// TODO Shapes implement both the Compiles interface and dbType (which
// is not exported but probably should be)
type Shape interface {
	aspect.Compiles
	Create(aspect.Dialect) (string, error)
}

type Point struct {
	X, Y float64
}

func (p Point) String() string {
	return fmt.Sprintf(`POINT(%f %f)`, p.X, p.Y)
}

// TODO parameterization
func (p Point) Compile(d aspect.Dialect, params *aspect.Parameters) (string, error) {
	return fmt.Sprintf(`ST_Point(%f, %f)`, p.X, p.Y), nil
}

func (p Point) Create(d aspect.Dialect) (string, error) {
	return "POINT", nil
}

type MultiPoint struct {
	Points []Point
}

// String returns the element's SQL using the default dialect.
func (p MultiPoint) String() string {
	compiled, _ := p.Compile(&postgres.PostGres{}, aspect.Params())
	return compiled
}

// TODO
func (p MultiPoint) Compile(d aspect.Dialect, params *aspect.Parameters) (string, error) {
	return "", nil
}

func (p MultiPoint) Create(d aspect.Dialect) (string, error) {
	return "MULTIPOINT", nil
}

type Linestring struct {
	Points []Point
}

// String returns the element's SQL using the default dialect.
func (p Linestring) String() string {
	compiled, _ := p.Compile(&postgres.PostGres{}, aspect.Params())
	return compiled
}

func (p Linestring) Compile(d aspect.Dialect, params *aspect.Parameters) (string, error) {
	return "", nil
}

func (p Linestring) Create(d aspect.Dialect) (string, error) {
	return "LINESTRING", nil
}

type Polygon struct {
	Exterior  Linestring
	Interiors []Linestring
}

// String returns the element's SQL using the default dialect.
func (p Polygon) String() string {
	compiled, _ := p.Compile(&postgres.PostGres{}, aspect.Params())
	return compiled
}

func (p Polygon) Compile(d aspect.Dialect, params *aspect.Parameters) (string, error) {
	return "", nil
}

func (p Polygon) Create(d aspect.Dialect) (string, error) {
	return "POLYGON", nil
}
