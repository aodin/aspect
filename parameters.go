package aspect

import ()

// Parameters holds a slice of interface{} parameters
type Parameters struct {
	args []interface{}
}

// Add adds a parameter to the parameter slice
func (p *Parameters) Add(i interface{}) int {
	p.args = append(p.args, i)
	return len(p.args)
}

func (p *Parameters) Args() []interface{} {
	return p.args
}

// Len returns the number of parameters in the slice
func (p *Parameters) Len() int {
	return len(p.args)
}

// Params creates a new Parameters instance
func Params() *Parameters {
	return &Parameters{}
}

// Parameter stores a single value of type interface{}
type Parameter struct {
	Value interface{}
}

func (p *Parameter) String() string {
	compiled, _ := p.Compile(&defaultDialect{}, Params())
	return compiled
}

// Parameter compilation is dialect dependent. For instance, dialects such
// as PostGres require the parameter index.
func (p *Parameter) Compile(d Dialect, params *Parameters) (string, error) {
	i := params.Add(p.Value)
	return d.Parameterize(i), nil
}
