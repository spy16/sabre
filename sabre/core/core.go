// Package core provides core contracts Sabre works with, builtin types of
// Sabre and a reader that can read builtins and be extended to read custom
// value types using reader macros.
package core

// Env represents the environment for LISP execution and maintains the value
// bindings created by the execution.
type Env interface {
	Eval(form Value) (Value, error)
	Bind(symbol string, v Value) error
	Resolve(symbol string) (Value, error)
	Parent() Env
}

// Value represents data/forms in sabre. This includes those emitted by Reader,
// values obtained as result of an evaluation etc..
type Value interface {
	// Eval should evaluate this value against the env and return
	// the resultant value or an evaluation error.
	Eval(env Env) (Value, error)

	// String should return the LISP representation of the value.
	String() string
}

// Invokable represents a value that can be invoked when it appears as the first
// entry in a list. For example, Keyword uses this to enable lookups in maps.
type Invokable interface {
	Value
	Invoke(env Env, args ...Value) (Value, error)
}

// Comparable can be implemented by Value types to support custom comparison logic.
// See Compare().
type Comparable interface {
	Value
	Compare(other Value) bool
}

// LazySeq implements a sequence of values that may be realized lazily.
type LazySeq interface {
	Value

	// First returns the first value of the sequence if not empty. Returns 'nil'
	// if empty.
	First() Value

	// Next returns a remaining sequence when the first value of the sequence
	// is excluded. 'nil' if the sequence is empty or has single item.
	Next() Seq

	// Cons returns a new sequence with the given value added as the first.
	Cons(v Value) Seq

	// Conj returns a new sequence which includes values from this sequence and
	// the arguments.
	Conj(vals ...Value) Seq
}

// Seq implementations represent a sequence of values such as List, Vector
// etc.
type Seq interface {
	LazySeq

	// Count returns the number of items in the map.
	Count() int
}

// Map represents any value that can store key-value pairs and provide fast
// lookups.
type Map interface {
	Value

	// Count returns the number of items in the map.
	Count() int

	// Keys returns a sequence of all keys in the map.
	Keys() Seq

	// Vals returns a sequence of all values in the map.
	Vals() Seq

	// HasKey checks if the map contains an entry with the given key.
	HasKey(key Value) bool

	// Get returns the value associated with the given key. Returns
	// ErrNotFound if the key not found.
	Get(key Value) (Value, error)

	// Assoc should return a new map which contains all the current values
	// with the given key-val pair added.
	Assoc(key, val Value) (Map, error)

	// Dissoc should return a new map which contains all the current entries
	// except the one with given key.
	Dissoc(key Value) (Map, error)
}

// Set represents a container for storing unique values.
type Set interface {
	Value

	Count() int
	Keys() Seq
	HasKey(key Value) bool
	Conj(vals ...Value) (Set, error)
}

// Attributable represents any value that supports dynamic attributes.
type Attributable interface {
	Value

	GetAttr(name string, defaultVal Value) Value
	SetAttr(name string, val Value) (Attributable, error)
}
