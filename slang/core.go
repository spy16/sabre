package slang

import (
	"github.com/spy16/sabre"
)

// Eval evaluates the first argument and returns the result.
func Eval(scope sabre.Scope, args []sabre.Value) (sabre.Value, error) {
	vals, err := evalValueList(scope, args)
	if err != nil {
		return nil, err
	}

	if err := verifyArgCount([]int{1}, args); err != nil {
		return nil, err
	}

	return vals[0].Eval(scope)
}

// Not returns the negated version of the argument value.
func Not(args []sabre.Value) (sabre.Value, error) {
	if err := verifyArgCount([]int{1}, args); err != nil {
		return nil, err
	}

	return sabre.Bool(!isTruthy(args[0])), nil
}
