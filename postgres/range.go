package postgres

import (
	"github.com/aodin/aspect"
)

type DateRange struct{}

var _ aspect.Type = DateRange{}

func (s DateRange) Create(d aspect.Dialect) (string, error) {
	return "daterange", nil
}

func (s DateRange) IsPrimaryKey() bool {
	return false
}

func (s DateRange) IsRequired() bool {
	return false
}

func (s DateRange) IsUnique() bool {
	return false
}

func (s DateRange) Validate(i interface{}) (interface{}, error) {
	return i, nil
}
