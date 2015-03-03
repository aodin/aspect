package aspect

// Type is the interface that all sql Types must implement
type Type interface {
	Creatable
	IsPrimaryKey() bool
	IsRequired() bool
	IsUnique() bool
	Validate(interface{}) (interface{}, error)
}
