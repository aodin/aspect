package aspect

import (
	"fmt"
	"reflect"
)

// fieldMap matches the given column names to the field names or `db` tags
// of the given struct. If not all columns were matched, the map is returned
// incomplete.
func fieldMap(columns []ColumnElem, i interface{}) (map[string]string, error) {
	// Get the type of the interface pointer
	t := reflect.TypeOf(i)
	if t.Kind() != reflect.Ptr {
		t = reflect.PtrTo(t)
	}

	// Map column name as key to field name as value
	alias := make(map[string]string)

	// There must be an underlying struct
	elem := t.Elem()
	if elem.Kind() != reflect.Struct {
		return alias, fmt.Errorf(
			"aspect: fieldMap only takes structs, received %s",
			elem.Kind(),
		)
	}

	for _, column := range columns {
		name := column.Name() // TODO full column name?
		// Duplicate column names will generate an error
		if _, exists := alias[name]; exists {
			return alias, fmt.Errorf(
				"aspect: found duplicate column name %s in fieldMap",
				name,
			)
		}

		// For each field, try the tag name, then the field name
		for i := 0; i < elem.NumField(); i += 1 {
			f := elem.Field(i)
			tag := f.Tag.Get("db")
			if tag == name || name == f.Name {
				alias[name] = f.Name
				break
			}
		}
	}
	return alias, nil
}