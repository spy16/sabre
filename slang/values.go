package slang

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spy16/sabre"
)

// TypeOf returns the type information object for the given argument.
func TypeOf(val sabre.Value) sabre.Value {
	return Type{rt: reflect.TypeOf(val)}
}

// IsType returns a Fn that checks if the value is of given type.
func IsType(rt reflect.Type) Fn {
	return func(vals []sabre.Value) (sabre.Value, error) {
		if err := verifyArgCount([]int{1}, vals); err != nil {
			return nil, err
		}

		target := reflect.TypeOf(vals[0])
		return sabre.Bool(target == rt), nil
	}
}

// MakeBool converts given argument to a boolean. Any truthy value
// is converted to true and else false.
func MakeBool(val sabre.Value) sabre.Bool {
	return sabre.Bool(IsTruthy(val))
}

// MakeInt converts given value to integer and returns.
func MakeInt(vals []sabre.Value) (sabre.Value, error) {
	if err := verifyArgCount([]int{1}, vals); err != nil {
		return nil, err
	}

	to := reflect.TypeOf(sabre.Int64(0))
	rv := reflect.ValueOf(vals[0])

	if !rv.Type().ConvertibleTo(to) {
		return nil, fmt.Errorf("cannot convert '%s' to '%s'", rv.Type(), to)
	}

	return rv.Convert(to).Interface().(sabre.Int64), nil
}

// MakeFloat converts given value to float and returns.
func MakeFloat(vals []sabre.Value) (sabre.Value, error) {
	if err := verifyArgCount([]int{1}, vals); err != nil {
		return nil, err
	}

	to := reflect.TypeOf(sabre.Float64(0))
	rv := reflect.ValueOf(vals[0])

	if !rv.Type().ConvertibleTo(to) {
		return nil, fmt.Errorf("cannot convert '%s' to '%s'", rv.Type(), to)
	}

	return rv.Convert(to).Interface().(sabre.Float64), nil
}

// MakeString returns stringified version of all args.
func MakeString(vals ...sabre.Value) sabre.Value {
	argc := len(vals)
	switch argc {
	case 0:
		return sabre.String("")

	case 1:
		nilVal := sabre.Nil{}
		if vals[0] == nilVal || vals[0] == nil {
			return sabre.String("")
		}

		return sabre.String(strings.Trim(vals[0].String(), "\""))

	default:
		var sb strings.Builder
		for _, v := range vals {
			sb.WriteString(strings.Trim(v.String(), "\""))
		}
		return sabre.String(sb.String())
	}
}

// makeContainer can make a composite type like list, set and vector from
// given args.
func makeContainer(targetType sabre.Value) Fn {
	return func(args []sabre.Value) (sabre.Value, error) {
		switch targetType.(type) {
		case *sabre.List:
			return &sabre.List{Values: args}, nil

		case sabre.Vector:
			return sabre.Vector{Values: args}, nil

		case sabre.Set:
			if err := verifyArgCount([]int{1}, args); err != nil {
				return nil, err
			}

			seq, ok := args[0].(sabre.Seq)
			if !ok {
				return nil, fmt.Errorf("can't create seq from '%s'",
					reflect.TypeOf(args[0]))
			}

			var seqVals sabre.Values
			for seq != nil && seq.First() != nil {
				seqVals = append(seqVals, seq.First())
				seq = seq.Next()
			}

			return sabre.Set{
				Values: sabre.Values(seqVals).Uniq(),
			}, nil
		}

		return nil, fmt.Errorf("cannot make container of type '%s'", reflect.TypeOf(targetType))
	}
}

// Type represents the type value of a given value. Type also implements
// Value type.
type Type struct {
	rt reflect.Type
}

// Eval returns the type value itself.
func (t Type) Eval(_ sabre.Scope) (sabre.Value, error) {
	return t, nil
}

func (t Type) String() string {
	return fmt.Sprintf("%v", t.rt)
}

// Invoke creates zero value of the given type.
func (t Type) Invoke(scope sabre.Scope, args ...sabre.Value) (sabre.Value, error) {
	return sabre.ValueOf(reflect.New(t.rt).Interface()), nil
}

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
