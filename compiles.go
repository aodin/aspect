package aspect

import ()

// The main SQL statement interface. All clauses must implement this
// interface in order to be executable.
type Compiles interface {
	String() string // For simple, parameter-less output
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
