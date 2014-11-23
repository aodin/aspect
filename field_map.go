package aspect

import (
	"fmt"
	"reflect"
)

type Field struct {
	Name      string
	OmitEmpty bool
}

// fieldMap matches the given column names to the field names or `db` tags
// of the given struct. If not all columns were matched, the map is returned
// incomplete.
func fieldMap(columns []ColumnElem, i interface{}) (map[string]Field, error) {
	// Get the type of the interface pointer
	t := reflect.TypeOf(i)
	if t.Kind() != reflect.Ptr {
		t = reflect.PtrTo(t)
	}

	// Map column name as key to field name as value
	alias := make(map[string]Field)

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
			tag, options := parseTag(f.Tag.Get("db"))

			// Skip the field
			if tag == "-" {
				continue
			}

			if tag == name || f.Name == name {
				alias[name] = Field{
					Name:      f.Name,
					OmitEmpty: options.Has("omitempty"),
				}
				break
			}
		}
	}
	return alias, nil
}

// Used in select statements - TODO should be generalized
func selectAlias(names []string, elem reflect.Type) (fields []int) {
	// TODO what about duplicate names?
	for _, name := range names {
		for i := 0; i < elem.NumField(); i += 1 {
			f := elem.Field(i)

			tag, _ := parseTag(f.Tag.Get("db"))

			// Skip the field
			if tag == "-" {
				continue
			}
			if tag == name || f.Name == name {
				fields = append(fields, i)
				break
			}
		}
	}
	return
}

// Used in select statements - TODO should be generalized
func selectIndex(names []string, elem reflect.Type) (fields []int) {
	for i := 0; i < elem.NumField(); i += 1 {
		fields = append(fields, i)
	}
	return
}

// TODO better way to pass columns than byusing the whole statement?
// TODO if this is better generalized then it can be used with UPDATE and
// DELETE statements.
func valuesMap(s InsertStmt, values Values) (map[string]Field, error) {
	alias := make(map[string]Field)
	for k, _ := range values {
		if !s.HasColumn(k) {
			return alias, fmt.Errorf(
				"aspect: cannot INSERT a value of key '%s' as it has no corresponding column",
				k,
			)
		}
		alias[k] = Field{Name: k}
	}
	return alias, nil
}
