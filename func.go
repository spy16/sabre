package sabre

import (
	"fmt"
	"reflect"
)

// SpecialFn implementations receive unevaluated list of forms during invoke.
// This allows such values to define their own evaluation rules.
type SpecialFn func(scope Scope, args []Value) (Value, error)

// Eval simply returns the special fn.
func (fn SpecialFn) Eval(_ Scope) (Value, error) { return fn, nil }

func (fn SpecialFn) String() string { return fmt.Sprintf("SpecialFn{%v}", reflect.ValueOf(fn)) }

// Invoke simply dispaatches the call to the wrapped function.
func (fn SpecialFn) Invoke(scope Scope, args []Value) (Value, error) { return fn(scope, args) }

type funcValue struct {
	rv reflect.Value
}

func (fn funcValue) Eval(_ Scope) (Value, error) { return fn, nil }

func (fn funcValue) String() string { return fmt.Sprintf("Fn{%s}", fn.rv) }

func (fn funcValue) Invoke(args []Value) (Value, error) {
	rt := fn.rv.Type()

	argVals, err := makeArgs(rt, args)
	if err != nil {
		return nil, err
	}

	retVals := fn.rv.Call(argVals)

	if rt.NumOut() == 0 {
		return nil, nil
	} else if rt.NumOut() == 1 {
		return ValueOf(retVals[0].Interface()), nil
	}

	var wrappedRetVals List
	for _, retVal := range retVals {
		wrappedRetVals = append(wrappedRetVals, ValueOf(retVal.Interface()))
	}

	return wrappedRetVals, nil
}

func makeArgs(rType reflect.Type, args []Value) ([]reflect.Value, error) {
	argVals := []reflect.Value{}

	if rType.IsVariadic() {
		nonVariadicLength := rType.NumIn() - 1
		for i := 0; i < nonVariadicLength; i++ {
			convertedArgVal, err := convertValueType(args[i], rType.In(i))
			if err != nil {
				return nil, err
			}

			argVals = append(argVals, convertedArgVal)
		}

		variadicType := rType.In(nonVariadicLength).Elem()
		for i := nonVariadicLength; i < len(args); i++ {
			convertedArgVal, err := convertValueType(args[i], variadicType)
			if err != nil {
				return nil, err
			}

			argVals = append(argVals, convertedArgVal)
		}

		return argVals, nil
	}

	if rType.NumIn() != len(args) {
		return nil, fmt.Errorf("call requires exactly %d arguments, got %d", rType.NumIn(), len(args))
	}

	for i := 0; i < rType.NumIn(); i++ {
		convertedArgVal, err := convertValueType(args[i], rType.In(i))
		if err != nil {
			return nil, err
		}

		argVals = append(argVals, convertedArgVal)
	}

	return argVals, nil
}

func convertValueType(v interface{}, expected reflect.Type) (reflect.Value, error) {
	val := newValue(v)
	if val.RVal.Type() == expected {
		return val.RVal, nil
	}

	converted, err := val.To(expected.Kind())
	if err != nil {
		if err == errConversion {
			return reflect.Value{}, fmt.Errorf("invalid argument type: expected=%s, actual=%s", expected, val.RVal.Type())
		}
		return reflect.Value{}, err
	}

	return reflect.ValueOf(converted), nil
}
