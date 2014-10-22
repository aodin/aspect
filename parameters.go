package aspect

import ()

// Parameters holds
type Parameters struct {
	// TODO Or just type Parameters []interface{} ?
	args []interface{}
}

func (p *Parameters) Add(i interface{}) int {
	p.args = append(p.args, i)
	return len(p.args)
}

func (p *Parameters) Args() []interface{} {
	return p.args
}

func (p *Parameters) Len() int {
	return len(p.args)
}

func Params() *Parameters {
	return &Parameters{args: make([]interface{}, 0)}
}

// Parameter stores a single value of type interface{}
type Parameter struct {
	// TODO Or just type Parameter []interface{}?
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
