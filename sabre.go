package sabre

import (
	"io"
	"strings"
)

// Eval evaluates the given form against the scope and returns the result
// of evaluation.
func Eval(scope Scope, form Value) (Value, error) {
	err := analyze(scope, form)
	if err != nil {
		return nil, err
	}

	return form.Eval(scope)
}

// ReadEval consumes data from reader 'r' till EOF, parses into forms
// and evaluates all the forms obtained and returns the result.
func ReadEval(scope Scope, r io.Reader) (Value, error) {
	mod, err := NewReader(r).All()
	if err != nil {
		return nil, err
	}

	return Eval(scope, mod)
}

// ReadEvalStr is a convenience wrapper for Eval that reads forms from
// string and evaluates for result.
func ReadEvalStr(scope Scope, src string) (Value, error) {
	return ReadEval(scope, strings.NewReader(src))
}

// Scope implementation is responsible for managing value bindings.
type Scope interface {
	Parent() Scope
	Bind(symbol string, v Value) error
	Resolve(symbol string) (Value, error)
}
