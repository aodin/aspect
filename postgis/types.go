package postgis

import (
	"fmt"
	"github.com/aodin/aspect"
)

// TODO They should take a shape instead

type GeometryPoint struct {
	Coord int
}

func (s GeometryPoint) Create(d aspect.Dialect) (string, error) {
	return fmt.Sprintf(`geometry(Point, %d)`, s.Coord), nil
}

type GeometryPolygon struct {
	Coord int
}

func (s GeometryPolygon) Create(d aspect.Dialect) (string, error) {
	return fmt.Sprintf(`geometry(Polygon, %d)`, s.Coord), nil
}
