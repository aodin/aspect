package aspect

import (
	"database/sql"
	"fmt"
	"reflect"
	"unicode"
)

const tagLabel = "db"

var scannerType = reflect.TypeOf((*sql.Scanner)(nil)).Elem()

type field struct {
	index   []int  // struct field indexes - with possible embedding
	column  string // SQL column name
	table   string // SQL table name
	options options
}

func (f field) Exists() bool {
	return len(f.index) > 0
}

func SelectFields(v interface{}) ([]field, error) {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Ptr {
		return nil, fmt.Errorf(
			"aspect: can only take select into pointer values, received %s",
			value.Kind(),
		)
	}
	return recurse([]int{}, reflect.TypeOf(v).Elem()), nil
}

func SelectFieldsFromElem(elem reflect.Type) []field {
	return recurse([]int{}, elem)
}

func recurse(indexes []int, elem reflect.Type) (fields []field) {
	if elem.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < elem.NumField(); i += 1 {
		f := elem.Field(i)
		if f.Type.Kind() == reflect.Struct && !reflect.PtrTo(f.Type).Implements(scannerType) {
			switch f.Type.String() { // TODO switch on the actual type
			case "time.Time":
			default:
				fields = append(fields, recurse(append(indexes, i), f.Type)...)
				continue
			}
		}

		// TODO ignore un-exported fields?
		// TODO fallback to json tags?

		// Field is valid, create a new Field
		var tag string
		field := field{index: append(indexes, i)}
		tag, field.options = parseTag(f.Tag.Get(tagLabel))
		if tag == "-" {
			continue
		} else if tag == "" {
			// Fallback to the column name, but only if exported
			if unicode.IsLower(rune(f.Name[0])) {
				continue
			}
			tag = f.Name
		}
		field.table, field.column = splitName(tag)
		fields = append(fields, field)
	}
	return
}

// AlignColumns will reorder the given fields array to match the columns.
// Columns that do not match fields will be given empty field structs.
func AlignColumns(columns []string, fields []field) []field {
	aligned := make([]field, len(columns))
	// TODO aliases? tables? check if the columns first matches the fields?
	for i, column := range columns {
		for _, field := range fields {
			if field.column == column {
				aligned[i] = field
				break
			}
		}
	}
	return aligned
}
