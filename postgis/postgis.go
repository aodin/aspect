package postgis

import (
	"fmt"
	"github.com/aodin/aspect"
)

type Point struct {
	Latitude  float64
	Longitude float64
}

func (p Point) Compile(d aspect.Dialect, params *aspect.Parameters) (string, error) {
	return fmt.Sprintf(`ST_GeometryFromText('POINT(%f %f)', 4326)`, p.Longitude, p.Latitude), nil
}

func Within(c aspect.ColumnElem, p Point) aspect.Clause {
	return aspect.MakeFuncClause(
		"ST_Within",
		aspect.MakeArrayClause(", ", []aspect.Clause{p, c}),
	)
}

func GeoJSON(c aspect.ColumnElem) aspect.ColumnElem {
	return c.SetInner(aspect.MakeFuncClause(
		"ST_AsGeoJSON",
		c.Inner(),
	))
}
