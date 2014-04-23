package aspect

import ()

// The main SQL statement interface. All clauses must implement this
// interface in order to be executable.
type Compiles interface {
	Compile(Dialect, *Parameters) (string, error)
}

// Perform compilation of the given statement with the given dialect.
// Ignore the parameters.
func CompileWith(c Compiles, d Dialect) (string, error) {
	return c.Compile(d, Params())
}

func CompileWithParams(c Compiles, d Dialect, p *Parameters) (string, error) {
	return c.Compile(d, p)
}

// TODO Or just type Parameters []interface{} ?
type Parameters struct {
	args []interface{}
}

func (p *Parameters) Add(i interface{}) int {
	p.args = append(p.args, i)
	return len(p.args)
}

func (p *Parameters) Len() int {
	return len(p.args)
}

func Params() *Parameters {
	return &Parameters{args: make([]interface{}, 0)}
}
