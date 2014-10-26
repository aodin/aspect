package aspect

import ()

// The main SQL statement interface. All clauses must implement this
// interface in order to be an executable statement or fragment.
type Compiles interface {
	Compile(Dialect, *Parameters) (string, error)
}
