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

func SelectFields(v interface{}) (fields []field, err error) {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Ptr {
		err = fmt.Errorf(
			"aspect: can only take select into pointer values, received %s",
			value.Kind(),
		)
		return
	}
	return recurse([]int{}, reflect.TypeOf(v).Elem()), nil
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
