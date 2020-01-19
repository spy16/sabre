package core

import (
	"fmt"
	"reflect"

	"github.com/spy16/sabre"
)

// Fn implements invokable with simple functions.
type Fn func(vals []sabre.Value) (sabre.Value, error)

// Eval simply returns the value.
func (fn Fn) Eval(_ sabre.Scope) (sabre.Value, error) {
	return fn, nil
}

func (fn Fn) String() string {
	return fmt.Sprintf("%s", reflect.ValueOf(fn).Type())
}

// Invoke evaluates all the args against the scope and dispatches the
// evaluated list as args to the wrapped function.
func (fn Fn) Invoke(scope sabre.Scope, args ...sabre.Value) (sabre.Value, error) {
	vals, err := evalValueList(scope, args)
	if err != nil {
		return nil, err
	}

	return fn(vals)
}

// SpecialFn implements Invokable using a function value.
type SpecialFn func(scope sabre.Scope, args []sabre.Value) (sabre.Value, error)

// Eval simply returns the special fn.
func (fn SpecialFn) Eval(_ sabre.Scope) (sabre.Value, error) {
	return fn, nil
}

func (fn SpecialFn) String() string {
	return fmt.Sprintf("%s", reflect.ValueOf(fn).Type())
}

// Invoke simply dispatches the call to the wrapped function.
func (fn SpecialFn) Invoke(scope sabre.Scope, args ...sabre.Value) (sabre.Value, error) {
	return fn(scope, args)
}
