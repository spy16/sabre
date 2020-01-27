package core

import (
	"fmt"
	"reflect"

	"github.com/spy16/sabre"
)

// TypeOf returns the type information object for the given argument.
func TypeOf(vals []sabre.Value) (sabre.Value, error) {
	if err := verifyArgCount([]int{1}, vals); err != nil {
		return nil, err
	}

	return Type{rt: reflect.TypeOf(vals[0])}, nil
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
func MakeBool(vals []sabre.Value) (sabre.Value, error) {
	if err := verifyArgCount([]int{1}, vals); err != nil {
		return nil, err
	}

	return sabre.Bool(isTruthy(vals[0])), nil
}

// MakeString returns stringified version of all args.
func MakeString(vals []sabre.Value) (sabre.Value, error) {
	return stringFromVals(vals), nil
}

// makeContainer can make a composite type like list, set and vector from
// given args.
func makeContainer(targetType sabre.Value) Fn {
	return func(vals []sabre.Value) (sabre.Value, error) {
		switch targetType.(type) {
		case sabre.List:
			return sabre.List{Values: vals}, nil

		case sabre.Vector:
			return sabre.Vector{Values: vals}, nil

		case sabre.Set:
			return sabre.Set{Items: vals}, nil
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
