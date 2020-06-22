package core

import (
	"fmt"
	"reflect"

	"github.com/spy16/sabre/runtime"
)

var (
	scopeType = reflect.TypeOf((*runtime.Runtime)(nil)).Elem()
	errorType = reflect.TypeOf((*error)(nil)).Elem()
)

// ValueOf converts a Go value to sabre Value. If 'v' is already a Value type, it
// is returned as is. Primitive Go values like string, rune, int, float, bool are
// converted to the right sabre Value types. Functions are converted to the wrapper
// 'Fn' type. Value of type 'reflect.Type' will be wrapped as 'Type' which enables
// initializing a value of that type when invoked. All other types will be wrapped
// using 'Any' type.
func ValueOf(v interface{}) runtime.Value {
	if v == nil {
		return runtime.Nil{}
	}

	if val, isValue := v.(runtime.Value); isValue {
		return val
	}

	if rt, ok := v.(reflect.Type); ok {
		return Type{T: rt}
	}

	rv := reflect.ValueOf(v)

	switch rv.Kind() {
	case reflect.Func:
		return reflectFn(rv)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return runtime.Int64(rv.Int())

	case reflect.Float32, reflect.Float64:
		return runtime.Float64(rv.Float())

	case reflect.String:
		return runtime.String(rv.String())

	case reflect.Uint8:
		return runtime.Char(rv.Uint())

	case reflect.Bool:
		return runtime.Bool(rv.Bool())

	default:
		// TODO: handle array & slice as list/vector.
		return Any{V: rv}
	}
}

// AccessMember accesses nested member field of the given object.
func AccessMember(obj runtime.Value, fields []string) (runtime.Value, error) {
	var v runtime.Value
	var err error

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

// Any can be used to wrap arbitrary Go value into Sabre scope.
type Any struct{ V reflect.Value }

// Eval returns itself.
func (any Any) Eval(_ runtime.Runtime) (runtime.Value, error) { return any, nil }

func (any Any) String() string { return fmt.Sprintf("Any{%v}", any.V) }

// Type represents the type value of a given value.
type Type struct{ T reflect.Type }

// Eval returns the type value itself.
func (t Type) Eval(_ runtime.Runtime) (runtime.Value, error) { return t, nil }

func (t Type) String() string { return fmt.Sprintf("%v", t.T) }

// Invoke creates zero value of the given type.
func (t Type) Invoke(scope runtime.Runtime, args ...runtime.Value) (runtime.Value, error) {
	if isKind(t.T, reflect.Interface, reflect.Chan, reflect.Func) {
		return nil, fmt.Errorf("type '%s' cannot be initialized", t.T)
	}

	argVals, err := runtime.EvalAll(scope, args)
	if err != nil {
		return nil, err
	}

	likeSeq := isKind(t.T, reflect.Slice, reflect.Array)
	if likeSeq {
		return runtime.NewSeq(argVals...), nil
	}

	return ValueOf(reflect.New(t.T).Elem().Interface()), nil
}

// reflectFn creates a wrapper Fn for the given Go function value using
// reflection.
func reflectFn(rv reflect.Value) runtime.Invokable {
	fw := wrapFunc(rv)
	return runtime.GoFunc(
		func(rt runtime.Runtime, args ...runtime.Value) (v runtime.Value, err error) {
			defer func() {
				if v := recover(); v != nil {
					err = fmt.Errorf("panic: %v", v)
				}
			}()

			args, err = runtime.EvalAll(rt, args)
			if err != nil {
				return nil, err
			}

			return fw.Call(rt, args...)
		})
}

func wrapFunc(rv reflect.Value) *funcWrapper {
	rt := rv.Type()

	minArgs := rt.NumIn()
	if rt.IsVariadic() {
		minArgs = minArgs - 1
	}

	passScope := (minArgs > 0) && (rt.In(0) == scopeType)
	lastOutIdx := rt.NumOut() - 1
	returnsErr := lastOutIdx >= 0 && rt.Out(lastOutIdx) == errorType
	if returnsErr {
		lastOutIdx-- // ignore error value from return values
	}

	return &funcWrapper{
		rv:         rv,
		rt:         rt,
		minArgs:    minArgs,
		passScope:  passScope,
		returnsErr: returnsErr,
		lastOutIdx: lastOutIdx,
	}
}

type funcWrapper struct {
	rv         reflect.Value
	rt         reflect.Type
	passScope  bool
	minArgs    int
	returnsErr bool
	lastOutIdx int
}

func (fw *funcWrapper) Call(scope runtime.Runtime, vals ...runtime.Value) (runtime.Value, error) {
	args := reflectValues(vals)
	if fw.passScope {
		args = append([]reflect.Value{reflect.ValueOf(scope)}, args...)
	}

	if err := fw.checkArgCount(len(args)); err != nil {
		return nil, err
	}

	args, err := fw.convertTypes(args...)
	if err != nil {
		return nil, err
	}

	return fw.wrapReturns(fw.rv.Call(args)...)
}

func (fw *funcWrapper) convertTypes(args ...reflect.Value) ([]reflect.Value, error) {
	var vals []reflect.Value

	for i := 0; i < fw.rt.NumIn(); i++ {
		if fw.rt.IsVariadic() && i == fw.rt.NumIn()-1 {
			c, err := convertArgsTo(fw.rt.In(i).Elem(), args[i:]...)
			if err != nil {
				return nil, err
			}
			vals = append(vals, c...)
			break
		}

		c, err := convertArgsTo(fw.rt.In(i), args[i])
		if err != nil {
			return nil, err
		}
		vals = append(vals, c...)
	}

	return vals, nil
}

func (fw *funcWrapper) checkArgCount(count int) error {
	if count != fw.minArgs {
		if fw.rt.IsVariadic() && count < fw.minArgs {
			return fmt.Errorf(
				"call requires at-least %d argument(s), got %d",
				fw.minArgs, count,
			)
		} else if !fw.rt.IsVariadic() && count > fw.minArgs {
			return fmt.Errorf(
				"call requires exactly %d argument(s), got %d",
				fw.minArgs, count,
			)
		}
	}

	return nil
}

func (fw *funcWrapper) wrapReturns(vals ...reflect.Value) (runtime.Value, error) {
	if fw.rt.NumOut() == 0 {
		return runtime.Nil{}, nil
	}

	if fw.returnsErr {
		errIndex := fw.lastOutIdx + 1
		if !vals[errIndex].IsNil() {
			return nil, vals[errIndex].Interface().(error)
		}

		if fw.rt.NumOut() == 1 {
			return runtime.Nil{}, nil
		}
	}

	wrapped := sabreValues(vals[0 : fw.lastOutIdx+1])
	if len(wrapped) == 1 {
		return wrapped[0], nil
	}

	return runtime.NewSeq(wrapped...), nil
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

func sabreValues(rvs []reflect.Value) []runtime.Value {
	var vals []runtime.Value
	for _, arg := range rvs {
		vals = append(vals, ValueOf(arg.Interface()))
	}
	return vals
}

func convertArgsTo(expected reflect.Type, args ...reflect.Value) ([]reflect.Value, error) {
	var converted []reflect.Value
	for _, arg := range args {
		actual := arg.Type()
		switch {
		case isAssignable(actual, expected):
			converted = append(converted, arg)

		case actual.ConvertibleTo(expected):
			converted = append(converted, arg.Convert(expected))

		default:
			return args, fmt.Errorf(
				"value of type '%s' cannot be converted to '%s'",
				actual, expected,
			)
		}
	}

	return converted, nil
}

func isAssignable(from, to reflect.Type) bool {
	return (from == to) || from.AssignableTo(to) ||
		(to.Kind() == reflect.Interface && from.Implements(to))
}

func reflectValues(args []runtime.Value) []reflect.Value {
	var rvs []reflect.Value
	for _, arg := range args {
		if any, ok := arg.(Any); ok {
			rvs = append(rvs, any.V)
		} else {
			rvs = append(rvs, reflect.ValueOf(arg))
		}
	}
	return rvs
}

func isKind(rt reflect.Type, kinds ...reflect.Kind) bool {
	for _, k := range kinds {
		if k == rt.Kind() {
			return true
		}
	}

	return false
}
