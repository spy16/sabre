package sabre

import (
	"fmt"
	"reflect"

	"github.com/spy16/sabre/sabre/runtime"
)

// ValueOf returns Sabre value for the given Go value.
func ValueOf(v interface{}) runtime.Value {
	return nil
}

// AccessMember accesses nested member field of the given object.
func AccessMember(obj runtime.Value, fields []string) (v runtime.Value, err error) {
	if len(fields) == 0 {
		return obj, nil
	}

	for len(fields) > 0 {
		v, err = accessOne(obj, fields[0])
		if err != nil {
			return nil, err
		}
		fields = fields[1:]
	}

	return v, nil
}

func accessOne(target runtime.Value, field string) (runtime.Value, error) {
	attr, ok := target.(runtime.Attributable)
	if ok {
		if val := attr.GetAttr(field, nil); val != nil {
			return val, nil
		}
	}

	rv := reflect.ValueOf(target)
	if rv.Type() == reflect.TypeOf(Any{}) {
		rv = rv.Interface().(Any).V
	}

	rVal, err := reflectAccess(rv, field)
	if err != nil {
		return nil, err
	}

	if isKind(rv.Type(), reflect.Chan, reflect.Array,
		reflect.Func, reflect.Ptr) && rv.IsNil() {
		return runtime.Nil{}, nil
	}

	return ValueOf(rVal.Interface()), nil
}

func reflectAccess(target reflect.Value, member string) (reflect.Value, error) {
	if member[0] >= 'a' && member[0] <= 'z' {
		return reflect.Value{}, fmt.Errorf("cannot access private member")
	}

	if _, found := target.Type().MethodByName(member); found {
		return target.MethodByName(member), nil
	}

	if target.Kind() == reflect.Ptr {
		target = target.Elem()
	}

	if _, found := target.Type().FieldByName(member); found {
		return target.FieldByName(member), nil
	}

	return reflect.Value{}, fmt.Errorf("value of type '%s' has no member named '%s'",
		target.Type(), member)
}

func isKind(rt reflect.Type, kinds ...reflect.Kind) bool {
	for _, k := range kinds {
		if k == rt.Kind() {
			return true
		}
	}
	return false
}

type Any struct {
	V reflect.Value
}
