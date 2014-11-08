package postgis

import (
	"fmt"

	"github.com/aodin/aspect"
)

type Geometry struct {
	Geom Shape
	SRID int
}

var _ aspect.Type = Geometry{}

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

func (g Geometry) IsPrimaryKey() bool {
	return false
}

func (g Geometry) IsRequired() bool {
	return false
}

func (g Geometry) IsUnique() bool {
	return false
}

func (g Geometry) Validate(i interface{}) (interface{}, error) {
	return i, nil
}
