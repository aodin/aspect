package aspect

import "sort"

// Values is a map of column names to parameters.
type Values map[string]interface{}

// Compile converts all key value pairs into a binary clauses.
// Since map iteration is non-deterministic, we'll sort the keys to
// produce repeatable SQL statements (especially for testing)
func (v Values) Compile(d Dialect, params *Parameters) (string, error) {
	clauses := make([]Clause, len(v))
	for i, key := range v.Keys() {
		clauses[i] = BinaryClause{
			Pre:  ColumnClause{name: key},
			Post: &Parameter{v[key]},
			Sep:  " = ",
		}
	}
	return ArrayClause{Clauses: clauses, Sep: ", "}.Compile(d, params)
}

// Keys returns the keys of the Values map in alphabetical order.
func (v Values) Keys() []string {
	keys := make([]string, len(v))
	var i int
	for key := range v {
		keys[i] = key
		i += 1
	}
	sort.Strings(keys)
	return keys
}
