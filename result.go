package aspect

import (
	"database/sql"
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

type Result struct {
	stmt string
	rows *sql.Rows
}

func (r *Result) Close() error {
	return r.rows.Close()
}

func (r *Result) Next() bool {
	return r.rows.Next()
}

// Return one result from the row as the given interface
func (r *Result) One(i interface{}) error {
	// Confirm that there is a row to return
	if ok := r.rows.Next(); !ok {
		return ErrNoResult
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
	} else if elem.Kind() == reflect.Map {
		// TODO This operation should not be repeated
		names, _ := r.rows.Columns()

		result := make([]interface{}, len(names))
		addresses := make([]interface{}, len(names))
		for index := range addresses {
			addresses[index] = &result[index]
		}
		if err := r.rows.Scan(addresses...); err != nil {
			return fmt.Errorf("aspect: error while scanning map: %s", err)
		}

		m, ok := elem.Interface().(map[string]interface{})
		if !ok {
			return ErrWrongMap
		}
		// Make a new mapping item
		// TODO Confirm that it is not a nil map
		// TODO what if there are redundant names?
		for index, name := range names {
			m[name] = result[index]
		}
		return nil
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
	} else if elem.Kind() == reflect.Map {
		names, _ := r.rows.Columns()
		result := make([]interface{}, len(names))
		addresses := make([]interface{}, len(names))
		for index := range addresses {
			addresses[index] = &result[index]
		}
		for r.rows.Next() {
			// Create a new slice element
			newElem := reflect.MakeMap(elem)

			if err := r.rows.Scan(addresses...); err != nil {
				return fmt.Errorf("aspect: error while scanning map: %s", err)
			}

			m, ok := newElem.Interface().(map[string]interface{})
			if !ok {
				return ErrWrongMap
			}
			// Make a new mapping item
			// TODO Confirm that it is not a nil map
			// TODO what if there are redundant names?
			for index, name := range names {
				m[name] = result[index]
			}
			sliceValue.Set(reflect.Append(sliceValue, newElem))
		}
		return nil
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
