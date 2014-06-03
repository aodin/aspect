package postgis

import (
	"fmt"
	"github.com/aodin/aspect"
)

type Geometry struct {
	Geom Shape
	SRID int
}

func (g Geometry) Create(d aspect.Dialect) (string, error) {
	inner, err := g.Geom.Create(d)
	if err != nil {
		return "", err
	}
	if g.SRID == 0 {
		return fmt.Sprintf(`geometry(%s)`, inner), nil
	}
	return fmt.Sprintf(`geometry(%s, %d)`, inner, g.SRID), nil
}
