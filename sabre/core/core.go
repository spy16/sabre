// Package core provides core contracts Sabre works with, builtin types of
// Sabre and a reader that can read builtins and be extended to read custom
// value types using reader macros.
package core

import (
	"errors"
	"fmt"
)

var (
	// ErrSkip can be returned by reader macro to indicate a no-op form which
	// should be discarded (e.g., Comments).
	ErrSkip = errors.New("skip expr")

	// ErrEOF is returned by reader when stream ends prematurely to indicate
	// that more data is needed to complete the current form.
	ErrEOF = errors.New("unexpected EOF")

	// ErrNotFound should be returned by an env implementation when a binding is
	// not found or by values that implement Associative when an entry is not
	// found.
	ErrNotFound = errors.New("not found")
)

// Env represents the environment for LISP execution and maintains the value
// bindings created by the execution.
type Env interface {
	Eval(form Value) (Value, error)

	// Bind binds the value to the symbol.
	Bind(symbol string, v Value) error

	// Resolve returns the value bound for the the given symbol.
	Resolve(symbol string) (Value, error)

	// Parent returns the parent of this environment.
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

	// Invoke is called when this value appears as first item in a list. Remaining
	// items of the list will be passed un-evaluated as arguments.
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

	// Count returns number of items in the set.
	Count() int

	// Keys returns the items/keys in the set.
	Keys() Seq

	// HasKey returns true if the key is present in the set.
	HasKey(key Value) bool

	// Conj returns a new set with the vals conjoined.
	Conj(vals ...Value) (Set, error)
}

// Attributable represents any value that supports dynamic attributes.
type Attributable interface {
	Value

	GetAttr(name string, defaultVal Value) Value
	SetAttr(name string, val Value) (Attributable, error)
}

// NewErr returns a sabre error object with given err as cause. If err is already
// a sabre Error, simply returns copy of it with given position attached.
func NewErr(isRead bool, pos Position, err error) Error {
	if ee, ok := err.(Error); ok {
		ee.Position = pos
		return ee
	} else if ee, ok := err.(*Error); ok && ee != nil {
		err := *ee
		err.Position = pos
		return err
	}

	return Error{
		Position:  pos,
		Cause:     err,
		IsReadErr: isRead,
	}
}

// Error represents errors during read or evaluation stages.
type Error struct {
	Position
	IsReadErr bool
	Message   string
	Cause     error
	Form      Value
}

// Unwrap returns the underlying cause of this error.
func (err Error) Unwrap() error { return err.Cause }

func (err Error) Error() string {
	if e, ok := err.Cause.(Error); ok {
		return e.Error()
	}

	if err.IsReadErr {
		return fmt.Sprintf(
			"syntax error in '%s' (Line %d Col %d): %v",
			err.File, err.Line, err.Column, err.Cause,
		)
	}

	return fmt.Sprintf("eval-error in '%s' (at line %d:%d): %v",
		err.File, err.Line, err.Column, err.Cause,
	)
}
