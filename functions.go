package aspect

import ()

func Avg(c ColumnElem) ColumnElem {
	c.inner = FuncClause{Inner: c.inner, F: "AVG"}
	return c
}

func Sum(c ColumnElem) ColumnElem {
	c.inner = FuncClause{Inner: c.inner, F: "SUM"}
	return c
}

func Count(c ColumnElem) ColumnElem {
	c.inner = FuncClause{Inner: c.inner, F: "COUNT"}
	return c
}

func DateOf(c ColumnElem) ColumnElem {
	c.inner = FuncClause{Inner: c.inner, F: "DATE"}
	return c
}

func Max(c ColumnElem) ColumnElem {
	c.inner = FuncClause{Inner: c.inner, F: "MAX"}
	return c
}

func Lower(c ColumnElem) ColumnElem {
	c.inner = FuncClause{Inner: c.inner, F: "LOWER"}
	return c
}

func DatePart(c ColumnElem, part string) ColumnElem {
	// Add the given date part as a parameter
	c.inner = FuncClause{
		Inner: ArrayClause{
			Clauses: []Clause{StringClause{Name: part}, c.inner},
			Sep:     ", ",
		},
		F: "DATE_PART",
	}
	return c
}
