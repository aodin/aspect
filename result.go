package aspect

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrNoResult = errors.New("aspect: no result to return")
	ErrWrongMap = errors.New(
		"aspect: maps must have key type string and value type interface{}",
	)
)

// Scanner is used for building mock result rows for testing
type Scanner interface {
	Close() error
	Columns() ([]string, error)
	Err() error
	Next() bool
	Scan(...interface{}) error
}

type Result struct {
	stmt string
	rows Scanner
}

func (r *Result) Close() error {
	return r.rows.Close()
}

func (r *Result) Next() bool {
	return r.rows.Next()
}

// One returns a single row from Result. The destination interface must be
// a pointer to a struct or a native type.
func (r *Result) One(arg interface{}) error {
	// Confirm that there is at least one row to return
	if ok := r.rows.Next(); !ok {
		return ErrNoResult
	}

	value := reflect.ValueOf(arg)
	if value.Kind() != reflect.Ptr {
		return fmt.Errorf(
			"aspect: received a non-pointer destination for result.One",
		)
	}

	columns, err := r.rows.Columns()
	if err != nil {
		return fmt.Errorf(
			"aspect: error returning columns from result: %s",
			err,
		)
	}

	// Get the value of the given interface
	elem := reflect.Indirect(value)

	switch elem.Kind() {
	case reflect.Struct:
		// Build the alias object from the columns and given elem
		alias := selectAlias(columns, reflect.TypeOf(arg).Elem())

		// The alias must be equal in length to the columns or if
		// the alias returned empty, the number of exported fields
		// must equal the columns.
		if len(alias) != len(columns) {
			alias = selectIndex(columns, reflect.TypeOf(arg).Elem())
		}

		// If the alias still does not match the number of columns, error
		if len(alias) != len(columns) {
			return fmt.Errorf("aspect: to select into a untagged struct without matching names, the struct's number of exported fields must match the order (and type) of the expected result")
		}

		// Get an interface for each field and save a pointer to it
		dest := make([]interface{}, len(alias))
		for i, fieldIndex := range alias {
			dest[i] = elem.Field(fieldIndex).Addr().Interface()
		}

		if err := r.rows.Scan(dest...); err != nil {
			return err
		}
	case reflect.Map:
		fallthrough
	default:
		// Attempt to scan directly into the elem
		return r.rows.Scan(elem.Addr().Interface())

		return fmt.Errorf(
			"aspect: unsupported destination element for result.One: %s",
			elem.Kind(),
		)
	}
	return r.rows.Err()
}

// All returns all result rows into the given interface, which must be a
// pointer to a slice of either structs, values, or a native type.
func (r *Result) All(arg interface{}) error {
	argVal := reflect.ValueOf(arg)
	if argVal.Kind() != reflect.Ptr {
		return fmt.Errorf(
			"aspect: received a non-pointer destination for result.All",
		)
	}

	argElem := argVal.Elem()
	if argElem.Kind() != reflect.Slice {
		return fmt.Errorf(
			"aspect: receive a non-slice destination for result.All",
		)
	}

	// Get the type of the slice element
	elem := argElem.Type().Elem()

	columns, err := r.rows.Columns()
	if err != nil {
		return fmt.Errorf(
			"aspect: error returning columns from result: %s",
			err,
		)
	}

	switch elem.Kind() {
	case reflect.Struct:
		// Build the alias object from the columns and given elem
		alias := selectAlias(columns, elem)

		// The alias must be equal in length to the columns or if
		// the alias returned empty, the number of exported fields
		// must equal the columns.
		if len(alias) != len(columns) {
			alias = selectIndex(columns, elem)
		}

		// If the alias still does not match the number of columns, error
		if len(alias) != len(columns) {
			return fmt.Errorf("aspect: to select into a untagged struct without matching names, the struct's number of exported fields must match the order (and type) of the expected result")
		}

		for r.rows.Next() {
			// Create a new slice element
			newElem := reflect.New(elem).Elem()

			// Get an interface for each field and save a pointer to it
			dest := make([]interface{}, len(alias))
			for i, fieldIndex := range alias {
				dest[i] = newElem.Field(fieldIndex).Addr().Interface()
			}

			if err := r.rows.Scan(dest...); err != nil {
				return err
			}
			argElem.Set(reflect.Append(argElem, newElem))
		}

	case reflect.Map:
		fallthrough
	default:
		return fmt.Errorf(
			"aspect: unsupported destination element for result.All: %s",
			elem.Kind(),
		)
	}

	// // Is the slice element a struct?
	// if elem.Kind() == reflect.Struct {
	// 	// Iterate through the struct fields to build the receiver
	// } else if elem.Kind() == reflect.Map {
	// 	names, _ := r.rows.Columns()
	// 	result := make([]interface{}, len(names))
	// 	addresses := make([]interface{}, len(names))
	// 	for index := range addresses {
	// 		addresses[index] = &result[index]
	// 	}
	// 	for r.rows.Next() {
	// 		// Create a new slice element
	// 		newElem := reflect.MakeMap(elem)

	// 		if err := r.rows.Scan(addresses...); err != nil {
	// 			return fmt.Errorf("aspect: error while scanning map: %s", err)
	// 		}

	// 		m, ok := newElem.Interface().(map[string]interface{})
	// 		if !ok {
	// 			return ErrWrongMap
	// 		}
	// 		// Make a new mapping item
	// 		// TODO Confirm that it is not a nil map
	// 		// TODO what if there are redundant names?
	// 		for index, name := range names {
	// 			m[name] = result[index]
	// 		}
	// 		argElem.Set(reflect.Append(argElem, newElem))
	// 	}
	// 	return nil
	// } else {
	// 	// Create a single interface for receiving
	// 	for r.rows.Next() {
	// 		// Create a new slice element
	// 		newElem := reflect.New(elem).Elem()
	// 		if err := r.rows.Scan(newElem.Addr().Interface()); err != nil {
	// 			return err
	// 		}
	// 		argElem.Set(reflect.Append(argElem, newElem))
	// 	}

	// }

	return r.rows.Err()
}
