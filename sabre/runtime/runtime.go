package runtime

import (
	"errors"
)

var (
	// ErrNotFound should be returned by a Runtime implementation when a binding
	// is not found or by values that implement Associative when an entry is not
	// found.
	ErrNotFound = errors.New("not found")

	// ErrNoEval can be returned by Runtime implementations to indicate that no eval
	// rule was found for the given form.
	ErrNoEval = errors.New("eval rule undefined")
)

// Runtime represents the environment for LISP execution. Runtime  defines the
// evaluation rules for forms and also is responsible for maintaining the value
// bindings created by the execution.
type Runtime interface {
	// Eval evaluates the form against runtime and returns the result of eval.
	// Evaluating might have side-effects on the runtime state (For example, a
	// def special-form will create a new binding in the global context).
	Eval(form Value) (Value, error)
}

// Value represents any value.
type Value interface {
	// String should return the textual representation of the value.
	String() string
}

// Invokable represents a value that can be invoked when it appears as the first
// entry in a list. For example, Keyword uses this to enable lookups in maps.
type Invokable interface {
	Value

	// Invoke is called when this value appears as first item in a list. Remaining
	// items of the list will be passed un-evaluated as arguments.
	Invoke(rt Runtime, args ...Value) (Value, error)
}

// Seq implements a sequence of values (e.g., List) that may be realized lazily.
type Seq interface {
	Value

	// First returns the first value of the sequence if not empty. Returns 'nil'
	// if empty.
	First() Value

	// Next returns the remaining sequence when the first value of the sequence
	// is excluded. 'nil' if the sequence is empty or has single item.
	Next() Seq

	// Conj returns a new sequence which includes values from this sequence and
	// the arguments. Position of conjoined values is not part of the contract.
	Conj(vals ...Value) Seq

	// Count returns the number of items in the map.
	Count() int
}

// Seqable is any value that can be converted to a sequence.
type Seqable interface {
	Value

	// Seq returns the implementing value as a sequence.
	Seq() Seq
}

// Vector represents a container for values that provides fast index lookups and
// iterations.
type Vector interface {
	Seqable

	// Count returns the number of elements in the vector.
	Count() int

	// EntryAt returns the item at given index. Returns error if the index
	// is out of range.
	EntryAt(index int) (Value, error)

	// Conj returns a new vector with items appended.
	Conj(items ...Value) Vector

	// Assoc returns a new vector with the value at given index updated.
	// Returns error if the index is out of range.
	Assoc(index int, val Value) (Vector, error)
}

// Map represents any value that can store key-value pairs and provide fast
// lookups.
type Map interface {
	Seqable

	// Keys returns all the keys in the map as a sequence.
	Keys() Seq

	// Vals returns all the values in the map as a sequence.
	Vals() Seq

	// HasKey returns true if the map contains the given key.
	HasKey(key Value) bool

	// EntryAt returns the value associated with the given key. Returns nil
	// if the key is not found or key is not hashable.
	EntryAt(key Value) Value

	// Assoc should return a new map which contains all the current values
	// with the given key-val pair added.
	Assoc(key, val Value) (Map, error)

	// Dissoc should return a new map which contains all the current entries
	// except the one with given key.
	Dissoc(key Value) (Map, error)
}

// Set represents a container for storing unique values.
type Set interface {
	Seqable

	// HasKey returns true if the key is present in the set.
	HasKey(key Value) bool

	// Conj returns a new set with the vals conjoined.
	Conj(vals ...Value) (Set, error)

	// Disj returns a new set with the vals dis-joined.
	Disj(vals ...Value) (Set, error)
}
