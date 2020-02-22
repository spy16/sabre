package slang

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/spy16/sabre"
)

// Throw converts args to strings and returns an error with all the strings
// joined.
func Throw(scope sabre.Scope, args ...sabre.Value) error {
	return errors.New(strings.Trim(MakeString(args...).String(), "\""))
}

// ApplySeq invokes fn with argument list formed by realizing the sequence.
func ApplySeq(scope sabre.Scope, fn sabre.Invokable, seq sabre.Seq) (sabre.Value, error) {
	return fn.Invoke(scope, Realize(seq).Values...)
}

// Concat concatenates s1 and s2 and returns a new sequence.
func Concat(s1, s2 sabre.Seq) sabre.Seq {
	vals := Realize(s1)
	vals.Values = append(vals.Values, Realize(s2).Values...)
	return vals
}

// Realize realizes a sequence by continuously calling First() and Next()
// until the sequence becomes nil.
func Realize(seq sabre.Seq) *sabre.List {
	var vals []sabre.Value

	for seq != nil {
		v := seq.First()
		if v == nil {
			break
		}
		vals = append(vals, v)
		seq = seq.Next()
	}

	return &sabre.List{Values: vals}
}

// TypeOf returns the type information object for the given argument.
func TypeOf(v interface{}) sabre.Value {
	return sabre.ValueOf(reflect.TypeOf(v))
}

// Implements checks if given value implements the interface represented
// by 't'. Returns error if 't' does not represent an interface type.
func Implements(v interface{}, t sabre.Type) (bool, error) {
	if t.R.Kind() == reflect.Ptr {
		t.R = t.R.Elem()
	}

	if t.R.Kind() != reflect.Interface {
		return false, fmt.Errorf("type '%s' is not an interface type", t)
	}

	return reflect.TypeOf(v).Implements(t.R), nil
}

// ToType attempts to convert given sabre value to target type. Returns
// error if conversion not possible.
func ToType(val sabre.Value, to sabre.Type) (sabre.Value, error) {
	rv := reflect.ValueOf(val)
	if rv.Type().ConvertibleTo(to.R) || rv.Type().AssignableTo(to.R) {
		return sabre.ValueOf(rv.Convert(to.R).Interface()), nil
	}

	return nil, fmt.Errorf("cannot convert '%s' to '%s'", rv.Type(), to.R)
}

// Assert implements (assert <expr> message?).
func Assert(scope sabre.Scope, args []sabre.Value) (sabre.Value, error) {
	if err := verifyArgCount([]int{1, 2}, args); err != nil {
		return nil, err
	}

	test, err := sabre.Eval(scope, args[0])
	if err != nil {
		return nil, err
	}

	if isTruthy(test) {
		return nil, nil
	}

	if len(args) == 1 {
		return nil, fmt.Errorf("assertion failed: '%s'", args[0])
	}

	msg, err := sabre.Eval(scope, args[1])
	if err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("%v", msg)
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
	return Realize(seq.Cons(v))
}

// Conj appends the second argument as last element in the first seq argument
// and returns.
func Conj(seq sabre.Seq, args ...sabre.Value) sabre.Value {
	return Realize(seq.Conj(args...))
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

	res, err := sabre.Eval(scope, args[0])
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
			res, err = sabre.Eval(scope, f)

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

func isTruthy(v sabre.Value) bool {
	if v == nil || v == (sabre.Nil{}) {
		return false
	}

	if b, ok := v.(sabre.Bool); ok {
		return bool(b)
	}

	return true
}
