package aspect

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Type interface {
	Creatable
	IsPrimaryKey() bool
	IsRequired() bool
	IsUnique() bool
	Validate(interface{}) (interface{}, error)
}

// Provide a nullable boolean for Boolean type default values
var (
	internalTrue  bool = true
	internalFalse bool = false
	True               = &internalTrue
	False              = &internalFalse
)

// Text represents TEXT column types.
type Text struct{}

var _ Type = Text{}

// Create returns the syntax need to create this column in CREATE statements.
func (s Text) Create(d Dialect) (string, error) {
	compiled := "TEXT"
	return compiled, nil
}

func (s Text) IsPrimaryKey() bool {
	return false
}

func (s Text) IsRequired() bool {
	return false
}

func (s Text) IsUnique() bool {
	return false
}

func (s Text) Validate(i interface{}) (interface{}, error) {
	if _, ok := i.(string); !ok {
		return i, fmt.Errorf("value is of non string type %T", i)
	}
	return i, nil
}

// String represents VARCHAR column types.
type String struct {
	Length     int
	NotNull    bool
	Unique     bool
	PrimaryKey bool
}

var _ Type = String{}

// Create returns the syntax need to create this column in CREATE statements.
func (s String) Create(d Dialect) (string, error) {
	compiled := "VARCHAR"
	attrs := make([]string, 0)
	if s.Length != 0 {
		compiled += fmt.Sprintf("(%d)", s.Length)
	}
	if s.PrimaryKey {
		attrs = append(attrs, "PRIMARY KEY")
	}
	if s.NotNull {
		attrs = append(attrs, "NOT NULL")
	}
	if s.Unique {
		attrs = append(attrs, "UNIQUE")
	}
	if len(attrs) > 0 {
		compiled += fmt.Sprintf(" %s", strings.Join(attrs, " "))
	}
	return compiled, nil
}

func (s String) IsPrimaryKey() bool {
	return s.PrimaryKey
}

func (s String) IsRequired() bool {
	return s.NotNull
}

func (s String) IsUnique() bool {
	return s.PrimaryKey || s.Unique
}

func (s String) Validate(i interface{}) (interface{}, error) {
	value, ok := i.(string)
	if !ok {
		return i, fmt.Errorf("value is of non string type %T", i)
	}
	if s.Length != 0 && len(value) > s.Length {
		return i, fmt.Errorf(
			"value is %d characters long, max length is %d",
			len(value),
			s.Length,
		)
	}
	return i, nil
}

// Integer represents INTEGER column types.
type Integer struct {
	NotNull    bool
	Unique     bool
	PrimaryKey bool
}

var _ Type = Integer{}

// Create returns the syntax need to create this column in CREATE statements.
func (s Integer) Create(d Dialect) (string, error) {
	compiled := "INTEGER"
	attrs := make([]string, 0)
	if s.PrimaryKey {
		attrs = append(attrs, "PRIMARY KEY")
	}
	if s.NotNull {
		attrs = append(attrs, "NOT NULL")
	}
	if s.Unique {
		attrs = append(attrs, "UNIQUE")
	}
	if len(attrs) > 0 {
		compiled += fmt.Sprintf(" %s", strings.Join(attrs, " "))
	}
	return compiled, nil
}

func (s Integer) IsPrimaryKey() bool {
	return s.PrimaryKey
}

func (s Integer) IsRequired() bool {
	return s.NotNull
}

func (s Integer) IsUnique() bool {
	return s.PrimaryKey || s.Unique
}

func (s Integer) Validate(i interface{}) (interface{}, error) {
	switch t := i.(type) {
	case string:
		v, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return i, err
		}
		i = v
	case float64:
		v := int64(t)
		if t != float64(v) {
			return i, fmt.Errorf(
				"value is numeric, but not a whole number: %f",
				t,
			)
		}
		i = v
	case int:
	case int64:
	default:
		return i, fmt.Errorf("value is non-numeric type %T", t)
	}
	return i, nil
}

// BigInt represents BIGINT column types.
type BigInt struct {
	NotNull    bool
	Unique     bool
	PrimaryKey bool
}

var _ Type = BigInt{}

// Create returns the syntax need to create this column in CREATE statements.
func (s BigInt) Create(d Dialect) (string, error) {
	compiled := "BIGINT"
	attrs := make([]string, 0)
	if s.PrimaryKey {
		attrs = append(attrs, "PRIMARY KEY")
	}
	if s.NotNull {
		attrs = append(attrs, "NOT NULL")
	}
	if s.Unique {
		attrs = append(attrs, "UNIQUE")
	}
	if len(attrs) > 0 {
		compiled += fmt.Sprintf(" %s", strings.Join(attrs, " "))
	}
	return compiled, nil
}

func (s BigInt) IsPrimaryKey() bool {
	return s.PrimaryKey
}

func (s BigInt) IsRequired() bool {
	return s.NotNull
}

func (s BigInt) IsUnique() bool {
	return s.PrimaryKey || s.Unique
}

func (s BigInt) Validate(i interface{}) (interface{}, error) {
	switch t := i.(type) {
	case string:
		v, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return i, err
		}
		i = v
	case float64:
		v := int64(t)
		if t != float64(v) {
			return i, fmt.Errorf(
				"value is numeric, but not a whole number: %f",
				t,
			)
		}
		i = v
	case int:
	case int64:
	default:
		return i, fmt.Errorf("value is non-numeric type %T", t)
	}
	return i, nil
}

// Timestamp represents TIMESTAMP column types.
// TODO take a time.Location for timezone options
type Timestamp struct {
	NotNull      bool
	PrimaryKey   bool
	WithTimezone bool
	Default      string // TODO Should this be a clause that compiles?
}

var _ Type = Timestamp{}

// Create returns the syntax need to create this column in CREATE statements.
func (s Timestamp) Create(d Dialect) (string, error) {
	compiled := "TIMESTAMP"
	if s.WithTimezone {
		compiled += " WITH TIME ZONE"
	}
	if s.NotNull {
		compiled += " NOT NULL"
	}
	if s.Default != "" {
		compiled += fmt.Sprintf(" DEFAULT (%s)", s.Default)
	}
	return compiled, nil
}

func (s Timestamp) IsPrimaryKey() bool {
	return s.PrimaryKey
}

func (s Timestamp) IsRequired() bool {
	return s.NotNull && s.Default == ""
}

func (s Timestamp) IsUnique() bool {
	return s.PrimaryKey
}

func (s Timestamp) Validate(i interface{}) (interface{}, error) {
	if _, ok := i.(time.Time); !ok {
		return i, fmt.Errorf("value is of non-time type %T", i)
	}
	return i, nil
}

// Date represents DATE column types.
type Date struct {
	NotNull    bool
	PrimaryKey bool
	Unique     bool
}

var _ Type = Date{}

// Create returns the syntax need to create this column in CREATE statements.
func (s Date) Create(d Dialect) (string, error) {
	compiled := "DATE"
	attrs := make([]string, 0)
	if s.PrimaryKey {
		attrs = append(attrs, "PRIMARY KEY")
	}
	if s.NotNull {
		attrs = append(attrs, "NOT NULL")
	}
	if s.Unique {
		attrs = append(attrs, "UNIQUE")
	}
	if len(attrs) > 0 {
		compiled += fmt.Sprintf(" %s", strings.Join(attrs, " "))
	}
	return compiled, nil
}

func (s Date) IsPrimaryKey() bool {
	return s.PrimaryKey
}

func (s Date) IsRequired() bool {
	return s.NotNull
}

func (s Date) IsUnique() bool {
	return s.PrimaryKey
}

func (s Date) Validate(i interface{}) (interface{}, error) {
	if _, ok := i.(time.Time); !ok {
		return i, fmt.Errorf("value is of non-time type %T", i)
	}
	return i, nil
}

// Boolean represents BOOL column types. It contains a Default field that
// can be left nil, or set with the included variables True and False.
type Boolean struct {
	NotNull bool
	Default *bool
}

var _ Type = Boolean{}

// Create returns the syntax need to create this column in CREATE statements.
func (s Boolean) Create(d Dialect) (string, error) {
	compiled := "BOOL"
	if s.NotNull {
		compiled += " NOT NULL"
	}
	if s.Default != nil {
		compiled += strings.ToUpper(fmt.Sprintf(" DEFAULT %t", *s.Default))
	}
	return compiled, nil
}

func (s Boolean) IsPrimaryKey() bool {
	return false
}

func (s Boolean) IsRequired() bool {
	return s.NotNull && s.Default == nil
}

func (s Boolean) IsUnique() bool {
	return false
}

func (s Boolean) Validate(i interface{}) (interface{}, error) {
	// TODO parse boolean strings
	if _, ok := i.(bool); !ok {
		return i, fmt.Errorf("value is of non-bool type %T", i)
	}
	return i, nil
}

// Double represents DOUBLE PRECISION column types.
type Double struct {
	NotNull    bool
	Unique     bool
	PrimaryKey bool
}

var _ Type = Double{}

// Create returns the syntax need to create this column in CREATE statements.
func (s Double) Create(d Dialect) (string, error) {
	compiled := "DOUBLE PRECISION"
	attrs := make([]string, 0)
	if s.PrimaryKey {
		attrs = append(attrs, "PRIMARY KEY")
	}
	if s.NotNull {
		attrs = append(attrs, "NOT NULL")
	}
	if s.Unique {
		attrs = append(attrs, "UNIQUE")
	}
	if len(attrs) > 0 {
		compiled += fmt.Sprintf(" %s", strings.Join(attrs, " "))
	}
	return compiled, nil
}

func (s Double) IsPrimaryKey() bool {
	return s.PrimaryKey
}

func (s Double) IsRequired() bool {
	return s.NotNull
}

func (s Double) IsUnique() bool {
	return s.PrimaryKey || s.Unique
}

func (s Double) Validate(i interface{}) (interface{}, error) {
	switch t := i.(type) {
	case string:
		v, err := strconv.ParseFloat(t, 64)
		if err != nil {
			return i, err
		}
		i = v
	case float32:
	case float64:
	case int:
		i = float64(t)
	case int64:
		i = float64(t)
	default:
		return i, fmt.Errorf("value is non-numeric type %T", t)
	}
	return i, nil
}

// Real represents REAL column types.
type Real struct {
	NotNull    bool
	Unique     bool
	PrimaryKey bool
}

var _ Type = Real{}

// Create returns the syntax need to create this column in CREATE statements.
func (s Real) Create(d Dialect) (string, error) {
	compiled := "REAL"
	attrs := make([]string, 0)
	if s.PrimaryKey {
		attrs = append(attrs, "PRIMARY KEY")
	}
	if s.NotNull {
		attrs = append(attrs, "NOT NULL")
	}
	if s.Unique {
		attrs = append(attrs, "UNIQUE")
	}
	if len(attrs) > 0 {
		compiled += fmt.Sprintf(" %s", strings.Join(attrs, " "))
	}
	return compiled, nil
}

func (s Real) IsPrimaryKey() bool {
	return s.PrimaryKey
}

func (s Real) IsRequired() bool {
	return s.NotNull
}

func (s Real) IsUnique() bool {
	return s.PrimaryKey || s.Unique
}

func (s Real) Validate(i interface{}) (interface{}, error) {
	switch t := i.(type) {
	case string:
		v, err := strconv.ParseFloat(t, 64)
		if err != nil {
			return i, err
		}
		i = v
	case float32:
	case float64:
	case int:
		i = float64(t)
	case int64:
		i = float64(t)
	default:
		return i, fmt.Errorf("value is non-numeric type %T", t)
	}
	return i, nil
}
