package aspect

import (
	"database/sql"
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

// Exists returns true if the field contains a valid recursive field index
func (f field) Exists() bool {
	return len(f.index) > 0
}

// HasOption returns true if the field's options contains the given option
func (f field) HasOption(option string) bool {
	return f.options.Has(option)
}

type fields []field

// Empty returns true if all of its fields do not exist
func (f fields) Empty() bool {
	for _, field := range f {
		if field.Exists() {
			return false
		}
	}
	return true
}

// HasColumn returns true if the given column exists in the fields
func (f fields) HasColumn(column string) bool {
	for _, field := range f {
		// TODO what about table name
		if field.column == column {
			return true
		}
	}
	return false
}

// SelectFields returns the ordered list of fields from the given interface.
func SelectFields(v interface{}) fields {
	return recurse([]int{}, reflect.TypeOf(v).Elem())
}

// SelectFieldsFromElem returns the ordered list of fields from the given
// reflect Type
func SelectFieldsFromElem(elem reflect.Type) fields {
	return recurse([]int{}, elem)
}

func recurse(indexes []int, elem reflect.Type) (fields fields) {
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

		// A new array will not actually be allocated during every
		// append because capacity is being increased by 2 - make sure to
		// perform a copy to allocate new memory
		indexesCopy := make([]int, len(indexes))
		copy(indexesCopy, indexes)
		field := field{index: append(indexesCopy, i)}
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
func AlignColumns(columns []string, fields []field) fields {
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
