package postgis

import (
	"github.com/aodin/aspect"
)

// Returns the area of the given column
// For "geometry" type area is in SRID units.
// For "geography" area is in square meters.
func AreaOf(c aspect.ColumnElem) aspect.ColumnElem {
	return c.SetInner(
		aspect.FuncClause{Inner: c.Inner(), F: "ST_Area"},
	)
}

// Returns the area of the surface if it is a polygon or multi-polygon.
// For "geometry" type area is in SRID units.
// For "geography" area is in square meters.
func Area(s Shape) aspect.Clause {
	return aspect.FuncClause{
		Inner: s,
		F:     "ST_Area",
	}
}

// Within returns true if the given shape is completely inside the given column
func Within(c aspect.ColumnElem, s Shape) aspect.Clause {
	return aspect.FuncClause{
		Inner: aspect.ArrayClause{Clauses: []aspect.Clause{s, c}, Sep: ", "},
		F:     "ST_Within",
	}
}

func DWithin(c aspect.ColumnElem, s Shape, d int) aspect.Clause {
	return aspect.FuncClause{
		Inner: aspect.ArrayClause{
			Clauses: []aspect.Clause{s, c, aspect.IntClause{D: d}},
			Sep:     ", ",
		},
		F: "ST_DWithin",
	}
}

// ST_Intersects — Returns TRUE if the Geometries/Geography "spatially
// intersect in 2D" - (share any portion of space) and FALSE if they don't
// (they are Disjoint). For geography -- tolerance is 0.00001 meters (so any
// points that close are considered to intersect)
func Intersects(c aspect.ColumnElem, s Shape) aspect.Clause {
	return aspect.FuncClause{
		Inner: aspect.ArrayClause{Clauses: []aspect.Clause{s, c}, Sep: ", "},
		F:     "ST_Intersects",
	}
}
