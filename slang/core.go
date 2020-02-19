package slang

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/spy16/sabre"
)

// Assert implements (assert <expr> message?).
func Assert(scope sabre.Scope, args []sabre.Value) (sabre.Value, error) {
	if err := verifyArgCount([]int{1, 2}, args); err != nil {
		return nil, err
	}

	test, err := args[0].Eval(scope)
	if err != nil {
		return nil, err
	}

	if IsTruthy(test) {
		return nil, nil
	}

	if len(args) == 1 {
		return nil, fmt.Errorf("assertion failed: '%s'", args[0])
	}

	msg, err := args[1].Eval(scope)
	if err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("%v", msg)
}

// IsEmpty returns true if the given value is nil or is an empty seq.
// If the value is non-nil and non-seq, returns error.
func IsEmpty(v sabre.Value) (bool, error) {
	if v == (sabre.Nil{}) || v == nil {
		return true, nil
	}

	seq, isSeq := v.(sabre.Seq)
	if !isSeq {
		return false, fmt.Errorf("value of type '%s' is not a sequence", reflect.TypeOf(v))
	}

	first := seq.First()
	return (first == (sabre.Nil{}) || first == nil), nil
}

// IsTruthy returns true if the given value is truthy. Boolean true,
// and all non-nil values are considered truthy.
func IsTruthy(v sabre.Value) bool {
	if v == nil {
		return false
	}

	var sabreNil = sabre.Nil{}
	if v == sabreNil {
		return false
	}

	if b, ok := v.(sabre.Bool); ok {
		return bool(b)
	}

	return true
}

// IsSeq returns true if the given value is a Seq.
func IsSeq(v sabre.Value) bool {
	_, isSeq := v.(sabre.Seq)
	return isSeq
}

// First returns the first value from the given Seq value.
func First(seq sabre.Seq) sabre.Value {
	return seq.First()
}

// Next returns the values after the first value as a list.
func Next(seq sabre.Seq) sabre.Value {
	n := seq.Next()
	if n == nil {
		return sabre.Nil{}
	}

	return n
}

// Cons inserts the first argument as first element in the second seq argument
// and returns.
func Cons(v sabre.Value, seq sabre.Seq) sabre.Value {
	return seq.Cons(v)
}

// Conj appends the second argument as last element in the first seq argument
// and returns.
func Conj(seq sabre.Seq, args ...sabre.Value) sabre.Value {
	return seq.Conj(args...)
}

// Not returns the negated version of the argument value.
func Not(val sabre.Value) sabre.Value {
	return sabre.Bool(!IsTruthy(val))
}

// ThreadFirst threads the expressions through forms by inserting result of
// eval as first argument to next expr.
func ThreadFirst(scope sabre.Scope, args []sabre.Value) (sabre.Value, error) {
	return threadCall(scope, args, false)
}

// ThreadLast threads the expressions through forms by inserting result of
// eval as last argument to next expr.
func ThreadLast(scope sabre.Scope, args []sabre.Value) (sabre.Value, error) {
	return threadCall(scope, args, true)
}

func threadCall(scope sabre.Scope, args []sabre.Value, last bool) (sabre.Value, error) {
	if len(args) == 0 {
		return nil, errors.New("at-least 1 argument required")
	}

	res, err := args[0].Eval(scope)
	if err != nil {
		return nil, err
	}

	for args = args[1:]; len(args) > 0; args = args[1:] {
		form := args[0]

		switch f := form.(type) {
		case *sabre.List:
			if last {
				f.Values = append(f.Values, res)
			} else {
				f.Values = append([]sabre.Value{f.Values[0], res}, f.Values[1:]...)
			}
			res, err = f.Eval(scope)

		case sabre.Invokable:
			res, err = f.Invoke(scope, res)

		default:
			return nil, fmt.Errorf("%s is not invokable", reflect.TypeOf(res))
		}

		if err != nil {
			return nil, err
		}
	}

	return res, nil
}
