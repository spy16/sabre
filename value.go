package sabre

import (
	"fmt"
	"reflect"
)

// ValueOf converts a Go value to sabre Value type. Functions will be
// converted to the Func type.
// TODO: handle array & slice as list/vector.
func ValueOf(v interface{}) Value {
	if val, isValue := v.(Value); isValue {
		return val
	}

	if v == nil {
		return List(nil)
	}

	rv := reflect.ValueOf(v)

	switch rv.Kind() {
	case reflect.Func:
		return funcValue{rv: rv}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return Int64(rv.Int())

	case reflect.Float32, reflect.Float64:
		return Float64(rv.Float())

	case reflect.String:
		return String(rv.String())

	case reflect.Uint8:
		return Character(rv.Uint())

	default:
		return anyValue{rv: rv}
	}
}

type anyValue struct {
	rv reflect.Value
}

func (any anyValue) Eval(_ Scope) (Value, error) {
	return any, nil
}

func (any anyValue) String() string {
	return fmt.Sprintf("Any{%s}", any.rv)
}
