package aspect

import ()

// Both ColumnStruct and OrderedColumns will implement the Orderable interface
type Orderable interface {
	Orderable() OrderedColumn
}

// From the PostGres documentation:
// http://www.postgresql.org/docs/9.2/static/sql-select.html#SQL-ORDERBY
// * If not specified, ASC is assumed by default.
// * the default behavior is NULLS LAST when ASC is specified or implied, and
// NULLS FIRST when DESC is specified
type OrderedColumn struct {
	inner                       ColumnStruct
	desc, nullsFirst, nullsLast bool
}

func (o OrderedColumn) String() string {
	compiled, _ := o.Compile(&PostGres{}, Params())
	return compiled
}

func (o OrderedColumn) Compile(d Dialect, params *Parameters) (string, error) {
	// Call the compilation method of the embeded column
	compiled, err := o.inner.Compile(d, params)
	if err != nil {
		return "", err
	}

	if o.desc {
		compiled += " DESC"
	}
	if o.nullsFirst || o.nullsLast {
		if o.nullsFirst {
			compiled += " NULLS FIRST"
		} else {
			compiled += " NULLS LAST"
		}
	}
	return compiled, nil
}

func (o OrderedColumn) Orderable() OrderedColumn {
	return o
}

func (o OrderedColumn) Asc() OrderedColumn {
	o.desc = false
	return o
}

func (o OrderedColumn) Desc() OrderedColumn {
	o.desc = true
	return o
}

func (o OrderedColumn) NullsFirst() OrderedColumn {
	o.nullsFirst = true
	o.nullsLast = false
	return o
}

func (o OrderedColumn) NullsLast() OrderedColumn {
	o.nullsFirst = false
	o.nullsLast = true
	return o
}
