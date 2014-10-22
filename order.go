package aspect

import ()

// Both ColumnElem and OrderedColumns will implement the Orderable interface
type Orderable interface {
	Orderable() OrderedColumn
}

// OrderedColumn represents a ColumnElem that will be used in an ORDER BY
// clause within SELECT statements. It provides additional sorting
// features, such as ASC, DESC, NULLS FIRST, and NULLS LAST.
// If not specified, ASC is assumed by default.
// In Postgres: the default behavior is NULLS LAST when ASC is specified or
// implied, and NULLS FIRST when DESC is specified
// http://www.postgresql.org/docs/9.2/static/sql-select.html#SQL-ORDERBY
type OrderedColumn struct {
	inner                       ColumnElem
	desc, nullsFirst, nullsLast bool
}

// OrderedColumn should implement the Orderable interface
var _ Orderable = OrderedColumn{}

func (o OrderedColumn) String() string {
	compiled, _ := o.Compile(&defaultDialect{}, Params())
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
