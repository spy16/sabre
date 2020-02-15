package sabre

import (
	"fmt"
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

	v, err := form.Eval(scope)
	if err != nil {
		if _, ok := err.(EvalError); ok {
			return v, err
		}

		return v, EvalError{
			Position: getPosition(form),
			Form:     v,
			Cause:    err,
		}
	}

	return v, nil
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

// EvalError represents error during evaluation.
type EvalError struct {
	Position
	Cause error
	Form  Value
}

// Unwrap returns the underlying cause of this error.
func (ee EvalError) Unwrap() error {
	return ee.Cause
}

func (ee EvalError) Error() string {
	return fmt.Sprintf("eval error in '%s' (Line %d, Column %d): %v",
		ee.File, ee.Line, ee.Column, ee.Cause,
	)
}
