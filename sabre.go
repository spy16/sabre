package sabre

import (
	"errors"
	"io"
	"strings"
)

var errNoEval = errors.New("no eval rule")

// ReadEval parses the entire string, evaluates all the forms obtained and
// returns the result.
func ReadEval(scope Scope, r io.Reader) (Value, error) {
	mod, err := NewReader(r).All()
	if err != nil {
		return nil, err
	}

	return mod.Eval(scope)
}

// ReadEvalStr is a convinience wrapper for ReadEval that reads string
// and evaluates for result.
func ReadEvalStr(scope Scope, src string) (Value, error) {
	return ReadEval(scope, strings.NewReader(src))
}

// Scope implementation is responsible for managing bindings
type Scope interface {
	Bind(name string, v Value) error
	Resolve(name string) (Value, error)
}

// Value represents a LISP value.
type Value interface {
	// Eval should evaluate this value against the scope and return
	// the resultant value or an evaluation error.
	Eval(scope Scope) (Value, error)

	// String should return the LISP representation of the value.
	String() string
}
