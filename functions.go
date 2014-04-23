package aspect

import ()

func Avg(c ColumnStruct) ColumnStruct {
	c.inner = FuncClause{clause: c.inner, f: "AVG"}
	return c
}

func Count(c ColumnStruct) ColumnStruct {
	c.inner = FuncClause{clause: c.inner, f: "COUNT"}
	return c
}

func Date(c ColumnStruct) ColumnStruct {
	c.inner = FuncClause{clause: c.inner, f: "DATE"}
	return c
}

func Max(c ColumnStruct) ColumnStruct {
	c.inner = FuncClause{clause: c.inner, f: "MAX"}
	return c
}

func DatePart(c ColumnStruct, part string) ColumnStruct {
	// Add the given date part as a parameter
	c.inner = FuncClause{
		clause: ArrayClause{
			clauses: []Clause{&Parameter{part}, c.inner},
			sep:     ", ",
		},
		f: "DATE_PART",
	}
	return c
}
