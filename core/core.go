// Package core provides core facilities of Sabre including types such
// as Keyword, Symbol, List, Vector etc. and the core interfaces such
// as Value, Env, Seq etc.
package core

import (
	"fmt"
	"reflect"
)

// Value represents data/forms in sabre. This includes those emitted by
// Reader, values obtained as result of an evaluation etc.
type Value interface {
	// Source should return the LISP representation of the value type.
	Source() string
}

// Expr represents values/forms that support self-evaluation. For such
// value types, env might dispatch the evaluation request to the Eval()
// method.
type Expr interface {
	Value

	// Eval should evaluate the form against the given env and return the
	// result of evaluation. Never use this method directly. It should be
	// called by the Env implementation.
	Eval(env Env) (Value, error)
}

// Invokable represents any value that supports invocation. Vector, Fn
// etc support invocation.
type Invokable interface {
	Value

	// Invoke is called when this value is present as first form in a list
	// and gets the current env and rest of the list as arguments.
	Invoke(env Env, args ...Value) (Value, error)
}

// Comparable can be implemented by Value types to support comparison.
// See Compare().
type Comparable interface {
	Value

	// Compare should compare the implementing value to the 'other' value.
	Compare(other Value) bool
}

// EvalAll evaluates all the values in the list against the env and returns
// the result of all evaluations.
func EvalAll(env Env, vals []Value) ([]Value, error) {
	var result []Value

	for _, arg := range vals {
		v, err := env.Eval(arg)
		if err != nil {
			return nil, err
		}
		result = append(result, v)
	}

	return result, nil
}

// Compare compares two values in an identity independent manner. If v1
// implements `Compare(Value) bool` method, the comparison is delegated
// to it as `v1.Compare(v2)`.
func Compare(v1, v2 Value) bool {
	if (v1 == nil && v2 == nil) ||
		(v1 == (Nil{}) && v2 == (Nil{})) {
		return true
	}

	if cmp, ok := v1.(Comparable); ok {
		return cmp.Compare(v2)
	}

	return reflect.DeepEqual(v1, v2)
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
