package slang

import (
	"reflect"

	"github.com/spy16/sabre"
)

// Eval evaluates the first argument and returns the result.
func Eval(scope sabre.Scope, arg sabre.Value) (sabre.Value, error) {
	return arg.Eval(scope)
}

// Not returns the negated version of the argument value.
func Not(val sabre.Value) sabre.Value {
	return sabre.Bool(!isTruthy(val))
}

// Equals compares 2 values using reflect.DeepEqual and returns
// the result.
func Equals(v1, v2 sabre.Value) bool {
	return reflect.DeepEqual(v1, v2)
}
