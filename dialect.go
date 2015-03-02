package aspect

import (
	"fmt"
	"log"
)

// Dialect is the common interface that all database drivers must implement.
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

// RegisterDialect adds the given Dialect to the registry at the given name.
func RegisterDialect(name string, d Dialect) {
	if d == nil {
		log.Panic("aspect: unable to register a nil Dialect")
	}
	if _, duplicate := dialects[name]; duplicate {
		log.Panic("aspect: a Dialect with this name already exists")
	}
	dialects[name] = d
}

// GetDialect returns the Dialect in the registry with the given name. An
// error will be returned if no Dialect with that name exists.
func GetDialect(name string) (Dialect, error) {
	d, ok := dialects[name]
	if !ok {
		return nil, fmt.Errorf(
			"aspect: unknown Dialect %s (did you remember to import it?)",
			name,
		)
	}
	return d, nil
}

// MustGetDialect returns the Dialect in the registry with the given name.
// It will panic if no Dialect with that name is found.
func MustGetDialect(name string) Dialect {
	dialect, err := GetDialect(name)
	if err != nil {
		log.Panic(err)
	}
	return dialect
}
