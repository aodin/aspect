package aspect

import ()

// ValidateInsert will check that the given struct or sql.Values is valid
// for insert on the given table.
func ValidateInsert(table *TableElem, v interface{}) *Errors {
	errors := EmptyErrors()

	// TODO Call fieldMap only if a struct
	fields, _ := fieldMap(table.Columns(), v)

	// Does the insert have all required fields?
	for _, column := range table.Columns() {
		if column.typ.IsRequired() {
			if _, ok := fields[column.name]; !ok {
				errors.SetField(column.name, "required")
			}
		}
	}
	if errors.Exist() {
		return errors
	}
	return nil
}
