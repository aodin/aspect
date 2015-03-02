package postgis

import (
	"github.com/aodin/aspect"
)

// Return the Well-Known Binary (WKB) representation of the geometry/geography
// without SRID meta data
func AsBinary(c aspect.ColumnElem) aspect.ColumnElem {
	return c.SetInner(
		aspect.FuncClause{Inner: c.Inner(), F: "ST_AsBinary"},
	)
}

// Return the Well-Known Text (WKT) representation of the geometry with
// SRID meta data
func AsEWKT(c aspect.ColumnElem) aspect.ColumnElem {
	return c.SetInner(
		aspect.FuncClause{Inner: c.Inner(), F: "ST_AsEWKT"},
	)
}

// Return the geometry as a GeoJSON element
func AsGeoJSON(c aspect.ColumnElem) aspect.ColumnElem {
	return c.SetInner(
		aspect.FuncClause{Inner: c.Inner(), F: "ST_AsGeoJSON"},
	)
}

// Return the geometry as a GML version 2 or 3 element
func AsGML(c aspect.ColumnElem) aspect.ColumnElem {
	return c.SetInner(
		aspect.FuncClause{Inner: c.Inner(), F: "ST_AsGML"},
	)
}

// Return the geometry as a KML element with default version=2 and precision=15
func AsKML(c aspect.ColumnElem) aspect.ColumnElem {
	return c.SetInner(
		aspect.FuncClause{Inner: c.Inner(), F: "ST_AsKML"},
	)
}

// Return the geometry as a KML element.
func AsKMLVersion(c aspect.ColumnElem, version, maxdigits int) aspect.ColumnElem {
	return c.SetInner(
		aspect.FuncClause{
			Inner: aspect.ArrayClause{
				Clauses: []aspect.Clause{
					aspect.IntClause{D: version},
					c.Inner(),
					aspect.IntClause{D: maxdigits},
				},
				Sep: ", ",
			},
			F: "ST_AsKML",
		},
	)
}

// Returns a Geometry in SVG path data given a geometry or geography object
func AsSVG(c aspect.ColumnElem) aspect.ColumnElem {
	return c.SetInner(
		aspect.FuncClause{Inner: c.Inner(), F: "ST_AsSVG"},
	)
}

// Return a GeoHash representation of the geometry.
func GeoHash(c aspect.ColumnElem) aspect.ColumnElem {
	return c.SetInner(
		aspect.FuncClause{Inner: c.Inner(), F: "ST_GeoHash"},
	)
}

// Return the Well-Known Text (WKT) representation of the geometry/geography
// without SRID metadata
func AsText(c aspect.ColumnElem) aspect.ColumnElem {
	return c.SetInner(
		aspect.FuncClause{Inner: c.Inner(), F: "ST_AsText"},
	)
}

// Return the Degrees, Minutes, Seconds representation of the given point
func AsLatLon(c aspect.ColumnElem) aspect.ColumnElem {
	return c.SetInner(
		aspect.FuncClause{Inner: c.Inner(), F: "ST_AsLatLonText"},
	)
}
