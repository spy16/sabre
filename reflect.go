package sabre

import (
	"fmt"
	"reflect"
)

// ValueOf converts a Go value to sabre Value type. Functions will be
// converted to the Func type. Other primitive Go types like string, rune,
// int (variants), float (variants) are converted to the right sabre Value
// types. If 'v' is already Value type, then it will be returned without
// conversion.
func ValueOf(v interface{}) Value {
	if val, isValue := v.(Value); isValue {
		return val
	}

	if v == nil {
		return Nil{}
	}

	rv := reflect.ValueOf(v)

	switch rv.Kind() {
	case reflect.Func:
		return reflectFn(rv)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return Int64(rv.Int())

	case reflect.Float32, reflect.Float64:
		return Float64(rv.Float())

	case reflect.String:
		return String(rv.String())

	case reflect.Uint8:
		return Character(rv.Uint())

	case reflect.Bool:
		return Bool(rv.Bool())

	default:
		// TODO: handle array & slice as list/vector.
		return anyValue{rv: rv}
	}
}

func reflectFn(rv reflect.Value) GoFunc {
	return func(scope Scope, args []Value) (_ Value, err error) {
		defer func() {
			if v := recover(); v != nil {
				if e, ok := v.(error); ok {
					err = e
				} else {
					err = fmt.Errorf("panic: %v", v)
				}
			}
		}()

		rt := rv.Type()
		argVals := reflectValues(args)

		if err := checkArgCount(rt, len(argVals)); err != nil {
			return nil, err
		}

		if err := checkArgTypes(rt, argVals); err != nil {
			return nil, err
		}

		retVals := rv.Call(argVals)

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
}

type anyValue struct{ rv reflect.Value }

func (any anyValue) Eval(_ Scope) (Value, error) { return any, nil }
func (any anyValue) String() string              { return fmt.Sprintf("Any{%v}", any.rv) }

func checkArgTypes(rt reflect.Type, args []reflect.Value) error {
	required := minArgs(rt)

	i := 0
	for ; i < required; i++ {
		if rt.In(i) != args[i].Type() {
			return fmt.Errorf("invalid argument type: expected=%s, actual=%s", rt.In(i), args[i].Type())
		}
	}

	if rt.IsVariadic() {
		expected := rt.In(i).Elem()
		for ; i < len(args); i++ {
			if expected != args[i].Type() {
				return fmt.Errorf("invalid argument type: expected=%s, actual=%s", expected, args[i].Type())
			}
		}
	}

	return nil
}

func reflectValues(args []Value) []reflect.Value {
	var rvs []reflect.Value

	for _, arg := range args {
		rvs = append(rvs, reflect.ValueOf(arg))
	}

	return rvs
}

func minArgs(rt reflect.Type) int {
	if rt.IsVariadic() {
		return rt.NumIn() - 1
	}

	return rt.NumIn()
}

func checkArgCount(rt reflect.Type, argCount int) error {
	if !rt.IsVariadic() {
		if argCount != minArgs(rt) {
			return fmt.Errorf("call requires exactly %d argument(s), got %d", rt.NumIn(), argCount)
		}
		return nil
	}

	required := minArgs(rt)
	if argCount < required {
		return fmt.Errorf("call requires at-least %d argument(s), got %d", required, argCount)
	}

	return nil
}
