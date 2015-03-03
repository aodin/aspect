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

	columns, err := r.rows.Columns()
	if err != nil {
		return fmt.Errorf(
			"aspect: error returning columns from result: %s",
			err,
		)
	}

	value := reflect.ValueOf(arg)
	if value.Kind() == reflect.Map {
		values, ok := arg.(Values)
		if !ok {
			return fmt.Errorf("aspect: maps as destinations are only allowed if they are of type aspect.Values")
		}

		// TODO scan directly into values?
		addr := make([]interface{}, len(columns))
		dest := make([]interface{}, len(columns))
		for i, _ := range addr {
			dest[i] = &addr[i]
		}

		if err := r.rows.Scan(dest...); err != nil {
			return fmt.Errorf("aspect: error while scanning map: %s", err)
		}

		for i, name := range columns {
			values[name] = addr[i]
		}
		return r.rows.Err()

	} else if value.Kind() != reflect.Ptr {
		return fmt.Errorf(
			"aspect: received a non-pointer destination for result.One",
		)
	}

	// Get the value of the given interface
	elem := reflect.Indirect(value)

	switch elem.Kind() {
	case reflect.Struct:
		// Build the fields of the given struct
		// TODO this operation could be cached
		fields, err := SelectFields(arg)
		if err != nil {
			return err
		}

		// Align the fields to the selected columns
		// This will discard unmatched fields
		// TODO struct mode? error if not all columns were matched?
		aligned := AlignColumns(columns, fields)

		// Get an interface for each field and save a pointer to it
		dest := make([]interface{}, len(aligned))
		for i, field := range aligned {
			// If the field does not exist, the value will be discarded
			if !field.Exists() {
				dest[i] = &dest[i]
				continue
			}

			// Recursively get an interface to the elem's fields
			var fieldElem reflect.Value = elem
			for _, index := range field.index {
				fieldElem = fieldElem.Field(index)
			}
			dest[i] = fieldElem.Addr().Interface()
		}

		if err := r.rows.Scan(dest...); err != nil {
			return fmt.Errorf("aspect: error while scanning struct: %s", err)
		}

	case reflect.Slice:
		return fmt.Errorf("aspect: cannot scan single results into slices")

	default:
		if len(columns) != 1 {
			return fmt.Errorf(
				"aspect: unsupported destination for multi-column result: %s",
				elem.Kind(),
			)
		}
		// Attempt to scan directly into the elem
		return r.rows.Scan(elem.Addr().Interface())
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

		// Build the fields of the given struct
		// TODO this operation could be cached
		fields := SelectFieldsFromElem(elem)

		// Align the fields to the selected columns
		// This will discard unmatched fields
		// TODO struct mode? error if not all columns were matched?
		aligned := AlignColumns(columns, fields)

		// Is there an existing slice element for this result?
		n := argElem.Len()

		// The number of results that hve been scanned
		var scanned int

		for r.rows.Next() {
			if scanned < n {
				// Scan into an existing element
				newElem := argElem.Index(scanned)

				// Get an interface for each field and save a pointer to it
				dest := make([]interface{}, len(aligned))
				for i, field := range aligned {
					// If the field does not exist, the value will be discarded
					if !field.Exists() {
						dest[i] = &dest[i]
						continue
					}

					// Recursively get an interface to the elem's fields
					var fieldElem reflect.Value = newElem
					for _, index := range field.index {
						fieldElem = fieldElem.Field(index)
					}
					dest[i] = fieldElem.Addr().Interface()
				}

				if err := r.rows.Scan(dest...); err != nil {
					return err
				}
			} else {
				// Create a new slice element
				newElem := reflect.New(elem).Elem()

				// Get an interface for each field and save a pointer to it
				dest := make([]interface{}, len(aligned))
				for i, field := range aligned {
					// If the field does not exist, the value will be discarded
					if !field.Exists() {
						dest[i] = &dest[i]
						continue
					}

					// Recursively get an interface to the elem's fields
					var fieldElem reflect.Value = newElem
					for _, index := range field.index {
						fieldElem = fieldElem.Field(index)
					}
					dest[i] = fieldElem.Addr().Interface()
				}

				if err := r.rows.Scan(dest...); err != nil {
					return err
				}
				argElem.Set(reflect.Append(argElem, newElem))
			}
			scanned += 1
		}

	case reflect.Map:
		_, ok := arg.(*[]Values)
		if !ok {
			return fmt.Errorf("aspect: slices of maps are only allowed if they are of type aspect.Values")
		}

		for r.rows.Next() {
			values := Values{}

			// TODO scan directly into values?
			addr := make([]interface{}, len(columns))
			dest := make([]interface{}, len(columns))
			for i, _ := range addr {
				dest[i] = &addr[i]
			}

			if err := r.rows.Scan(dest...); err != nil {
				return fmt.Errorf("aspect: error while scanning map: %s", err)
			}

			for i, name := range columns {
				values[name] = addr[i]
			}

			argElem.Set(reflect.Append(argElem, reflect.ValueOf(values)))
		}

	default:
		// Single column results can be scanned into native types
		if len(columns) != 1 {
			return fmt.Errorf(
				"aspect: unsupported destination for multi-column result: %s",
				elem.Kind(),
			)
		}
		for r.rows.Next() {
			newElem := reflect.New(elem).Elem()
			if err := r.rows.Scan(newElem.Addr().Interface()); err != nil {
				return err
			}
			argElem.Set(reflect.Append(argElem, newElem))
		}
	}

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

	return r.rows.Err()
}
