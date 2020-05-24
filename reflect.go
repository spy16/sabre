package sabre

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spy16/sabre/core"
)

var (
	scopeType = reflect.TypeOf((*core.Env)(nil)).Elem()
	errorType = reflect.TypeOf((*error)(nil)).Elem()
)

// ValueOf converts a Go value to sabre Value type. If 'v' is already a Value
// type, it is returned as is. Primitive Go values like string, rune, int, float,
// bool are converted to the right sabre Value types. Functions are converted to
// the wrapper 'Fn' type. Value of type 'reflect.Type' will be wrapped as 'Type'
// which enables initializing a value of that type when invoked. All other types
// will be wrapped using 'Any' type.
func ValueOf(v interface{}) core.Value {
	if v == nil {
		return core.Nil{}
	}

	if val, isValue := v.(core.Value); isValue {
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
		return core.Int64(rv.Int())

	case reflect.Float32, reflect.Float64:
		return core.Float64(rv.Float())

	case reflect.String:
		return core.String(rv.String())

	case reflect.Uint8:
		return core.Character(rv.Uint())

	case reflect.Bool:
		return core.Bool(rv.Bool())

	default:
		// TODO: handle array & slice as list/vector.
		return Any{V: rv}
	}
}

// Any can be used to wrap arbitrary Go value into Sabre scope.
type Any struct{ V reflect.Value }

// Source returns a string version of the value.
func (any Any) Source() string { return fmt.Sprintf("Any{%v}", any.V) }

// Type represents the type value of a given value. Type also implements
// Value and Invokable. Invoking type creates zero value of the wrapped
// type.
type Type struct{ T reflect.Type }

// Source returns the string version of the value.
func (t Type) Source() string { return fmt.Sprintf("%v", t.T) }

// Invoke creates zero value of the given type.
func (t Type) Invoke(env core.Env, args ...core.Value) (core.Value, error) {
	if isKind(t.T, reflect.Interface, reflect.Chan, reflect.Func) {
		return nil, fmt.Errorf("type '%s' cannot be initialized", t.T)
	}

	argVals, err := core.EvalAll(env, args)
	if err != nil {
		return nil, err
	}

	switch t.T {
	case reflect.TypeOf((*core.List)(nil)):
		return &core.List{Values: argVals}, nil

	case reflect.TypeOf(core.Vector{}):
		return core.Vector{Values: argVals}, nil

	case reflect.TypeOf(core.Set{}):
		return core.Set{Values: core.Values(argVals).Uniq()}, nil
	}

	likeSeq := isKind(t.T, reflect.Slice, reflect.Array)
	if likeSeq {
		return core.Values(argVals), nil
	}

	return ValueOf(reflect.New(t.T).Elem().Interface()), nil
}

// reflectFn creates a wrapper Fn for the given Go function value using
// reflection.
func reflectFn(rv reflect.Value) *core.Fn {
	fw := wrapFunc(rv)
	return &core.Fn{
		Args:     fw.argNames(),
		Variadic: rv.Type().IsVariadic(),
		Func: func(env core.Env, args []core.Value) (_ core.Value, err error) {
			defer func() {
				if v := recover(); v != nil {
					err = fmt.Errorf("panic: %v", v)
				}
			}()

			args, err = core.EvalAll(env, args)
			if err != nil {
				return nil, err
			}

			return fw.Call(env, args...)
		},
	}
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

func (fw *funcWrapper) Call(env core.Env, vals ...core.Value) (core.Value, error) {
	args := reflectValues(vals)
	if fw.passScope {
		args = append([]reflect.Value{reflect.ValueOf(env)}, args...)
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

func (fw *funcWrapper) argNames() []string {
	cleanArgName := func(t reflect.Type) string {
		return strings.Replace(t.String(), "sabre.", "", -1)
	}

	var argNames []string

	i := 0
	for ; i < fw.minArgs; i++ {
		argNames = append(argNames, cleanArgName(fw.rt.In(i)))
	}

	if fw.rt.IsVariadic() {
		argNames = append(argNames, cleanArgName(fw.rt.In(i).Elem()))
	}

	return argNames
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

func (fw *funcWrapper) wrapReturns(vals ...reflect.Value) (core.Value, error) {
	if fw.rt.NumOut() == 0 {
		return core.Nil{}, nil
	}

	if fw.returnsErr {
		errIndex := fw.lastOutIdx + 1
		if !vals[errIndex].IsNil() {
			return nil, vals[errIndex].Interface().(error)
		}

		if fw.rt.NumOut() == 1 {
			return core.Nil{}, nil
		}
	}

	wrapped := sabreValues(vals[0 : fw.lastOutIdx+1])
	if len(wrapped) == 1 {
		return wrapped[0], nil
	}

	return core.Values(wrapped), nil
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

func reflectValues(args []core.Value) []reflect.Value {
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

func sabreValues(rvs []reflect.Value) []core.Value {
	var vals []core.Value
	for _, arg := range rvs {
		vals = append(vals, ValueOf(arg.Interface()))
	}
	return vals
}

func isKind(rt reflect.Type, kinds ...reflect.Kind) bool {
	for _, k := range kinds {
		if k == rt.Kind() {
			return true
		}
	}

	return false
}
