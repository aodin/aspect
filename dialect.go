package aspect

import (
	"fmt"
)

type Dialect interface {
	Parameterize(int) string
}

// Test dialect - uses postgres style parameterization
type defaultDialect struct{}

func (d *defaultDialect) Parameterize(i int) string {
	return fmt.Sprintf(`$%d`, i)
}

// Registry of available dialects
var dialects = make(map[string]Dialect)

func RegisterDialect(name string, d Dialect) {
	if d == nil {
		panic("aspect: unable to register a nil Dialect")
	}
	if _, duplicate := dialects[name]; duplicate {
		panic("aspect: a Dialect with this name already exists")
	}
	dialects[name] = d
}

func GetDialect(name string) (Dialect, error) {
	d, ok := dialects[name]
	if !ok {
		return nil, fmt.Errorf("aspect: unknown Dialect %s (did you remember to import it?)", name)
	}
	return d, nil
}
