// Package runtime provides runtime contracts Sabre works with, builtin types of
// Sabre and a reader that can read builtins and be extended to read custom value
// types using reader macros. All primitive values implemented are immutable.
package runtime

import (
	"errors"
	"fmt"
)

// ErrNotFound should be returned by an env implementation when a binding is not
// found or by values that implement Associative when an entry is not found.
var ErrNotFound = errors.New("not found")

// Runtime represents the environment for LISP execution and maintains the value
// bindings created by the execution.
type Runtime interface {
	// Eval evaluates the form in this runtime. Runtime might customize the eval
	// rules for different values, but in most cases, eval will be dispatched to
	// to Eval() method of value.
	Eval(form Value) (Value, error)

	// Bind binds the value to the symbol. Returns error if the symbol contains
	// invalid character or the bind fails for some other reasons.
	Bind(symbol string, v Value) error

	// Resolve returns the value bound for the the given symbol. Resolve returns
	// ErrNotFound if the symbol has no binding.
	Resolve(symbol string) (Value, error)

	// Parent returns the parent of this environment. If returned value is nil,
	// it is the root Runtime.
	Parent() Runtime
}

// Value represents data/forms in sabre. This includes those emitted by Reader,
// values obtained as result of an evaluation etc..
type Value interface {
	// Eval should evaluate this value against the runtime and return the
	// resultant value or an evaluation error.
	Eval(rt Runtime) (Value, error)

	// String should return the LISP representation of the value.
	String() string
}

// Hashable represents any value that is hashable and can be used as key in hash
// based collection types (HashMap, HashSet etc.).
type Hashable interface {
	Value

	// Hash should return a byte sequence that incorporates the contents of the
	// value that should be considered in hashing. Hashing implementations will
	// use the returned byte sequence to generate the hashcode.
	Hash() []byte
}

// Invokable represents a value that can be invoked when it appears as the first
// entry in a list. For example, Keyword uses this to enable lookups in maps.
type Invokable interface {
	Value

	// Invoke is called when this value appears as first item in a list. Remaining
	// items of the list will be passed un-evaluated as arguments.
	Invoke(env Runtime, args ...Value) (Value, error)
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

	// Cons returns a new sequence with the given value added as the first.
	Cons(v Value) Seq

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
	// if the key is not found or is not hashable.
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
	Conj(vals ...Value) Set

	// Disj returns a new set with the vals dis-joined.
	Disj(vals ...Value) Set
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

// Position represents the positional information about a value read
// by reader.
type Position struct {
	File   string
	Line   int
	Column int
}

// GetPos returns the file, line and column values.
func (pi Position) GetPos() (file string, line, col int) {
	return pi.File, pi.Line, pi.Column
}

// SetPos sets the position information.
func (pi *Position) SetPos(file string, line, col int) {
	pi.File = file
	pi.Line = line
	pi.Column = col
}

func (pi Position) String() string {
	if pi.File == "" {
		pi.File = "<unknown>"
	}

	return fmt.Sprintf("%s:%d:%d", pi.File, pi.Line, pi.Column)
}
