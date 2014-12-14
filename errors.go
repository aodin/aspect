package aspect

import "fmt"

type Errors struct {
	Meta   []error
	Fields map[string]error
}

func (errors *Errors) Exist() bool {
	return (len(errors.Meta) > 0 || len(errors.Fields) > 0)
}

func (errors *Errors) SetField(field, msg string, args ...interface{}) {
	errors.Fields[field] = fmt.Errorf(msg, args...)
}

func EmptyErrors() *Errors {
	return &Errors{
		Fields: make(map[string]error),
	}
}
