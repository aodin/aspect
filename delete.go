package aspect

import (
	"fmt"
	"reflect"
)

// DeleteStmt is the internal representation of a DELETE statement.
type DeleteStmt struct {
	table   *TableElem
	args    []interface{}
	cond    Clause
	err     error
	pkField int
}

// String outputs the parameter-less DELETE statement in a neutral dialect.
func (stmt DeleteStmt) String() string {
	compiled, _ := stmt.Compile(&defaultDialect{}, Params())
	return compiled
}

// Compile outputs the DELETE statement using the given dialect and parameters.
// An error may be returned because of a pre-existing error or because
// an error occurred during compilation.
func (stmt DeleteStmt) Compile(d Dialect, params *Parameters) (string, error) {
	if stmt.err != nil {
		return "", stmt.err
	}
	compiled := fmt.Sprintf(`DELETE FROM "%s"`, stmt.table.Name)

	if stmt.cond != nil {
		cc, err := stmt.cond.Compile(d, params)
		if err != nil {
			return "", err
		}
		compiled += fmt.Sprintf(" WHERE %s", cc)
	}
	return compiled, nil
}

// Values sets the conditional clause using the given struct or
// slice of structs. The table must have a single primary key field.
func (stmt DeleteStmt) Values(arg interface{}) DeleteStmt {
	// The must be a primary key
	if len(stmt.table.pk) == 0 {
		stmt.err = fmt.Errorf("aspect: a table must have a primary key to delete by values")
		return stmt
	}

	if len(stmt.table.pk) > 1 {
		stmt.err = fmt.Errorf("aspect: deletion by composite primary keys is currently not supported")
		return stmt
	}

	// Get the actual pk column
	pk := stmt.table.C[stmt.table.pk[0]]

	// Determine the primary key field in the given struct (or slice)
	elem := reflect.Indirect(reflect.ValueOf(arg))

	switch elem.Kind() {
	case reflect.Struct:
		// Match the primary key column to the field
		if stmt.err = stmt.setPkField(arg, pk.Name()); stmt.err != nil {
			return stmt
		}

		// Set the conditional
		stmt.cond = pk.Equals(elem.Field(stmt.pkField).Interface())

	case reflect.Slice:
		if elem.Len() < 1 {
			stmt.err = fmt.Errorf(
				"aspect: values cannot be set for deletion by empty slices",
			)
			return stmt
		}

		// Only slices of structs are acceptable
		elem0 := elem.Index(0)
		if elem0.Kind() != reflect.Struct {
			stmt.err = fmt.Errorf(
				"aspect: unsupported slice element for deletion by values: %s",
				elem0.Kind(),
			)
			return stmt
		}

		// Build the pk parameters slice and In conditional
		var pks []interface{}
		for i := 0; i < elem.Len(); i++ {
			pks = append(pks, elem.Index(i).Field(stmt.pkField).Interface())
		}
		stmt.cond = pk.In(pks)

	default:
		stmt.err = fmt.Errorf(
			"aspect: unsupported type for deletion by values: %s",
			elem.Kind(),
		)
	}

	return stmt
}

// Where adds a conditional WHERE clause to the DELETE statement.
func (stmt DeleteStmt) Where(conds ...Clause) DeleteStmt {
	if len(conds) > 1 {
		// By default, multiple where clauses will be joined will AllOf
		stmt.cond = AllOf(conds...)
	} else if len(conds) == 1 {
		stmt.cond = conds[0]
	}
	return stmt
}

// Find the index of the field with a name or db tag that matches the given
// name.
func (stmt *DeleteStmt) setPkField(arg interface{}, name string) error {
	// Get the type of the interface pointer
	t := reflect.TypeOf(arg)
	if t.Kind() != reflect.Ptr {
		t = reflect.PtrTo(t)
	}

	// TODO There must be an underlying struct
	elem := t.Elem()

	// For each field, try the tag name, then the field name
	// TODO copied from fieldMap, generalize
	for i := 0; i < elem.NumField(); i += 1 {
		f := elem.Field(i)
		tag := f.Tag.Get("db")
		if tag == name || f.Name == name {
			stmt.pkField = i
			return nil
		}
	}
	return fmt.Errorf("aspect: could not find a field for column %s", name)
}

// Delete creates a DELETE statement for the given table.
func Delete(table *TableElem) (stmt DeleteStmt) {
	if table == nil {
		stmt.err = fmt.Errorf(
			"aspect: attempting to DELETE a nil table",
		)
		return
	}
	stmt.table = table
	return
}
