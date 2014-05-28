package aspect

import (
	"fmt"
	"reflect"
)

type DeleteStmt struct {
	table *TableElem
	args  []interface{}
	cond  Clause
}

func (stmt DeleteStmt) String() string {
	compiled, _ := stmt.Compile(&PostGres{}, Params())
	return compiled
}

func (stmt DeleteStmt) Compile(d Dialect, params *Parameters) (string, error) {
	compiled := fmt.Sprintf(`DELETE FROM "%s"`, stmt.table.Name)

	// TODO Add any existing arguments to the parameters

	if stmt.cond != nil {
		cc, err := stmt.cond.Compile(d, params)
		if err != nil {
			return "", err
		}
		compiled += fmt.Sprintf(" WHERE %s", cc)
	}
	return compiled, nil
}

func (stmt DeleteStmt) Where(cond Clause) DeleteStmt {
	stmt.cond = cond
	return stmt
}

// TODO Merge with insert/fieldAlias()
// Return the index of the field with the given name or db tag
func fieldIndex(i interface{}, name string) (int, error) {
	t := reflect.TypeOf(i)
	if t.Kind() != reflect.Ptr {
		t = reflect.PtrTo(t)
	}
	elem := t.Elem()
	if elem.Kind() != reflect.Struct {
		return 0, fmt.Errorf("Cannot take a field index of a non-struct")
	}
	for i := 0; i < elem.NumField(); i += 1 {
		f := elem.Field(i)
		if f.Name == name {
			return i, nil
		}
		tag := f.Tag.Get("db")
		if tag == name {
			return i, nil
		}
	}
	return 0, fmt.Errorf("No field matching this tag found")
}

func getFieldByIndex(i interface{}, index int) interface{} {
	elem := reflect.Indirect(reflect.ValueOf(i))
	if elem.Kind() != reflect.Struct {
		return nil
	}
	if index >= elem.NumField() {
		return nil
	}
	return elem.Field(index).Interface()
}

func Delete(table *TableElem, args ...interface{}) DeleteStmt {
	stmt := DeleteStmt{table: table}
	// If the table has a primary key and was given args, create a conditional
	// statement using its pk columns and the values from the given args
	if table.pk == nil || len(table.pk) == 0 || len(args) == 0 {
		return stmt
	}

	// TODO Create delayed errors or fail silently?
	// The args must be structs with either fields or fields with tags
	// matching the declared primary keys
	if len(table.pk) == 1 {
		// Does the given arg have a field matching the primary key?
		index, err := fieldIndex(args[0], table.pk[0])
		if err != nil {
			return stmt
		}

		// Get the ColumnElem named by the pk
		pk := table.C[table.pk[0]]

		// Turn the arguments primary key into a parameter
		// TODO All args must be the same - how to enforce?
		if len(args) > 1 {
			pks := make([]interface{}, len(args))
			for i, arg := range args {
				pks[i] = getFieldByIndex(arg, index)
			}
			stmt.cond = pk.In(pks)
		} else {
			stmt.cond = pk.Equals(getFieldByIndex(args[0], index))
		}
	} else {
		// TODO composite keys
	}

	return stmt
}
