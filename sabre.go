package sabre

import (
	"io"
	"strings"
)

// Eval consumes data from reader 'r' till EOF, parses into forms and
// evaluates all the forms obtained and returns the result.
func Eval(scope Scope, r io.Reader) (Value, error) {
	mod, err := NewReader(r).All()
	if err != nil {
		return nil, err
	}

	return mod.Eval(scope)
}

// EvalStr is a convenience wrapper for Eval that reads forms from string
// and evaluates for result.
func EvalStr(scope Scope, src string) (Value, error) {
	return Eval(scope, strings.NewReader(src))
}

// Scope implementation is responsible for managing value bindings.
type Scope interface {
	Parent() Scope
	Bind(symbol string, v Value) error
	Resolve(symbol string) (Value, error)
}
