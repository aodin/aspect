package aspect

// The main SQL statement interface. All clauses must implement this
// interface in order to be an executable statement or fragment.
type Compiles interface {
	String() string // Output a neutral dialect logging string
	Compile(Dialect, *Parameters) (string, error)
}
