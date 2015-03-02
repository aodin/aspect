package aspect

import (
	"sort"

	"github.com/stretchr/testify/assert"
)

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

// Diff returns the values in v that differ from the values in other.
// ISO 31-11: v \ other
func (v Values) Diff(other Values) Values {
	diff := Values{}
	for key, value := range v {
		if !assert.ObjectsAreEqual(value, other[key]) {
			diff[key] = value
		}
	}
	return diff
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

func (v Values) String() string {
	compiled, _ := v.Compile(&defaultDialect{}, Params())
	return compiled
}
