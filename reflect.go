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

type anyValue struct{ rv reflect.Value }

func (any anyValue) Eval(_ Scope) (Value, error) {
	return any, nil
}

func (any anyValue) String() string {
	return fmt.Sprintf("Any{%v}", any.rv)
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

		args, err = evalValueList(scope, args)
		if err != nil {
			return nil, err
		}

		rt := rv.Type()
		argVals := reflectValues(args)

		if minArgs(rt) > 0 {
			scopeRV := reflect.ValueOf(scope)
			if scopeRV.Type().AssignableTo(rt.In(0)) {
				argVals = append([]reflect.Value{scopeRV}, argVals...)
			}
		}

		if err := checkArgCount(rt, len(argVals)); err != nil {
			return nil, err
		}

		converted, err := convertArgTypes(rt, argVals)
		if err != nil {
			return nil, err
		}

		retVals := rv.Call(converted)
		return wrapReturnValues(rt, retVals)
	}
}

func wrapReturnValues(fn reflect.Type, vals []reflect.Value) (Value, error) {
	if fn.NumOut() == 0 {
		return Nil{}, nil
	}

	lastArgIdx := fn.NumOut() - 1
	isLastArgErr := fn.Out(lastArgIdx).Name() == "error"

	if isLastArgErr {
		if !vals[lastArgIdx].IsNil() {
			return nil, vals[lastArgIdx].Interface().(error)
		}

		if fn.NumOut() == 1 {
			return Nil{}, nil
		}
	}

	lastValIdx := lastArgIdx + 1
	if isLastArgErr {
		lastValIdx = lastValIdx - 1
	}

	var wrappedRetVals []Value
	for _, retVal := range vals[0:lastValIdx] {
		wrappedRetVals = append(wrappedRetVals, ValueOf(retVal.Interface()))
	}

	if len(wrappedRetVals) == 1 {
		return wrappedRetVals[0], nil
	}

	return &List{Values: wrappedRetVals}, nil
}

func convertArgTypes(rt reflect.Type, args []reflect.Value) ([]reflect.Value, error) {
	required := minArgs(rt)
	var converted []reflect.Value

	i := 0
	for ; i < required; i++ {
		c, err := convertArgType(rt.In(i), args[i])
		if err != nil {
			return nil, err
		}
		converted = append(converted, c)
	}

	if rt.IsVariadic() {
		expected := rt.In(i).Elem()
		for ; i < len(args); i++ {
			c, err := convertArgType(expected, args[i])
			if err != nil {
				return nil, err
			}
			converted = append(converted, c)
		}
	}

	return converted, nil
}

func convertArgType(expected reflect.Type, value reflect.Value) (reflect.Value, error) {
	actual := value.Type()
	switch {
	case actual == expected || actual.AssignableTo(expected):
		return value, nil

	case (expected.Kind() == reflect.Interface) &&
		actual.Implements(expected):
		return value, nil

	case actual.ConvertibleTo(expected):
		return value.Convert(expected), nil

	}

	return value, fmt.Errorf(
		"invalid argument type: expected=%s, actual=%s",
		expected, actual,
	)
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
			return fmt.Errorf(
				"call requires exactly %d argument(s), got %d",
				rt.NumIn(), argCount,
			)
		}
		return nil
	}

	required := minArgs(rt)
	if argCount < required {
		return fmt.Errorf(
			"call requires at-least %d argument(s), got %d",
			required, argCount,
		)
	}

	return nil
}
