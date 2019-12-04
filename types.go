package sabre

import (
	"errors"
	"reflect"
)

var errConversion = errors.New("cannot be converted")

func newValue(v interface{}) reflectVal {
	return reflectVal{
		RVal: reflect.ValueOf(v),
	}
}

type reflectVal struct {
	RVal reflect.Value
}

// To converts the value to requested kind if possible.
func (val *reflectVal) To(kind reflect.Kind) (interface{}, error) {
	switch kind {
	case reflect.Int, reflect.Int64:
		return val.ToInt64()
	case reflect.Float64:
		return val.ToFloat64()
	case reflect.String:
		return val.ToString()
	case reflect.Bool:
		return val.ToBool()
	case reflect.Interface:
		return val.RVal.Interface(), nil
	default:
		return nil, errConversion
	}
}

// ToInt64 attempts converting the value to int64.
func (val *reflectVal) ToInt64() (int64, error) {
	if val.isInt() {
		return val.RVal.Int(), nil
	} else if val.isFloat() {
		return int64(val.RVal.Float()), nil
	}

	return 0, errConversion
}

// ToFloat64 attempts converting the value to float64.
func (val *reflectVal) ToFloat64() (float64, error) {
	if val.isFloat() {
		return val.RVal.Float(), nil
	} else if val.isInt() {
		return float64(val.RVal.Int()), nil
	}

	return 0, errConversion
}

// ToBool attempts converting the value to bool.
func (val *reflectVal) ToBool() (bool, error) {
	if isKind(val.RVal, reflect.Bool) {
		return val.RVal.Bool(), nil
	}

	return false, errConversion
}

// ToString attempts converting the value to bool.
func (val *reflectVal) ToString() (string, error) {
	if isKind(val.RVal, reflect.String) {
		return val.RVal.String(), nil
	}

	return "", errConversion
}

func (val *reflectVal) isInt() bool {
	return isKind(val.RVal, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64)
}

func (val *reflectVal) isFloat() bool {
	return isKind(val.RVal, reflect.Float32, reflect.Float64)
}

func isKind(rval reflect.Value, kinds ...reflect.Kind) bool {
	for _, kind := range kinds {
		if rval.Kind() == kind {
			return true
		}
	}
	return false
}
