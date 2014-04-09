package aspect

import (
	"database/sql"
	"fmt"
	"reflect"
)

type Result struct {
	stmt string
	rows *sql.Rows
}

func (r *Result) Close() error {
	return r.rows.Close()
}

// Return one result from the row as the given interface
func (r *Result) One(i interface{}) error {
	// Confirm that there is a row to return
	if ok := r.rows.Next(); !ok {
		return fmt.Errorf("No more rows to return")
	}

	// Get the value of the given interface
	value := reflect.ValueOf(i)
	elem := reflect.Indirect(value)

	// Is the slice element a struct?
	if elem.Kind() == reflect.Struct {
		// TODO Cached destination construction?
		dest := make([]interface{}, elem.NumField())

		// Get an interface for each field and save a pointer to it
		for i := 0; i < elem.NumField(); i += 1 {
			dest[i] = elem.Field(i).Addr().Interface()
		}
		return r.rows.Scan(dest...)
	} else {
		return r.rows.Scan(elem.Addr().Interface())
	}
}

func (r *Result) All(i interface{}) error {
	// TODO Error checking
	// The destination interface must of type slice
	sliceValue := reflect.ValueOf(i).Elem()

	// Get the type of the slice element
	elem := sliceValue.Type().Elem()

	// Is the slice element a struct?
	if elem.Kind() == reflect.Struct {
		// Iterate through the struct fields to build the receiver
		for r.rows.Next() {
			// Create a new slice element
			newElem := reflect.New(elem).Elem()

			// TODO Cached destination construction?
			dest := make([]interface{}, newElem.NumField())

			// Get an interface for each field and save a pointer to it
			for i := 0; i < newElem.NumField(); i += 1 {
				dest[i] = newElem.Field(i).Addr().Interface()
			}
			if err := r.rows.Scan(dest...); err != nil {
				return err
			}
			sliceValue.Set(reflect.Append(sliceValue, newElem))
		}
	} else {
		// Create a single interface for receiving
		for r.rows.Next() {
			// Create a new slice element
			newElem := reflect.New(elem).Elem()
			if err := r.rows.Scan(newElem.Addr().Interface()); err != nil {
				return err
			}
			sliceValue.Set(reflect.Append(sliceValue, newElem))
		}

	}
	// TODO Duplication of scan err check?
	return r.rows.Err()
}
