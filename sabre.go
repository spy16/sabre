package sabre

// Scope implementations are responsible for managing bindings.
type Scope interface {
	Parent() Scope
	Get(name string) (Value, error)
}

// Value represents any LISP value that evaluate to itself.
type Value interface{}
