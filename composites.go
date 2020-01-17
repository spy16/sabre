package sabre

import (
	"fmt"
	"reflect"
)

// List represents an list of forms/vals. Evaluating a list leads to a
// function invocation.
type List []Value

// Eval performs an invocation.
func (lf List) Eval(scope Scope) (Value, error) {
	if len(lf) == 0 {
		return List(nil), nil
	}

	if isQuote(lf[0]) {
		return quoteValue(lf[1:])
	}

	target, err := lf[0].Eval(scope)
	if err != nil {
		return nil, err
	}

	fn, ok := target.(Invokable)
	if !ok {
		return nil, fmt.Errorf("cannot invoke value of type '%s'", reflect.TypeOf(target))
	}

	return fn.Invoke(scope, lf[1:]...)
}

func (lf List) String() string {
	if len(lf) == 2 && isQuote(lf[0]) {
		return fmt.Sprintf("'%s", lf[1])
	}

	return containerString(lf, "(", ")", " ")
}

// Vector represents a list of values. Unlike List type, evaluation of
// vector does not lead to function invoke.
type Vector []Value

// Eval evaluates each value in the vector form and returns the resultant
// values as new vector.
func (vf Vector) Eval(scope Scope) (Value, error) {
	vals, err := evalValueList(scope, vf)
	if err != nil {
		return nil, err
	}

	return Vector(vals), nil
}

// Invoke of a vector performs a index lookup. Only arity 1 is allowed
// and should be an integer value to be used as index.
func (vf Vector) Invoke(scope Scope, args ...Value) (Value, error) {
	if len(args) != 1 {
		return nil, arityErr(1, len(args), "")
	}

	v, err := args[0].Eval(scope)
	if err != nil {
		return nil, err
	}

	if !isInt(v) {
		return nil, fmt.Errorf("key must be integer")
	}

	index := reflect.ValueOf(v).Int()

	if int(index) >= len(vf) {
		return nil, fmt.Errorf("index out of bounds")
	}

	return vf[index], nil
}

func (vf Vector) String() string {
	return containerString(vf, "[", "]", " ")
}

// Module represents a group of forms. Evaluating a module leads to evaluation
// of each form in order and result will be the result of last evaluation.
type Module []Value

// Eval evaluates all the vals in the module body and returns the result of the
// last evaluation.
func (mod Module) Eval(scope Scope) (Value, error) {
	res, err := evalValueList(scope, mod)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return List(nil), nil
	}

	return res[len(res)-1], nil
}

func (mod Module) String() string { return containerString(mod, "", "\n", "\n") }

func evalValueList(scope Scope, vals []Value) ([]Value, error) {
	var result []Value

	for _, arg := range vals {
		v, err := arg.Eval(scope)
		if err != nil {
			return nil, err
		}

		result = append(result, v)
	}

	return result, nil
}

func isQuote(v Value) bool {
	sym, isSymbol := v.(Symbol)
	if !isSymbol {
		return false
	}

	return sym == "quote"
}

func quoteValue(args []Value) (Value, error) {
	if len(args) != 1 {
		return nil, arityErr(1, len(args), "")
	}

	return args[0], nil
}

func isInt(v interface{}) bool {
	return isKind(reflect.ValueOf(v),
		reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64,
	)
}

func isKind(rval reflect.Value, kinds ...reflect.Kind) bool {
	for _, kind := range kinds {
		if rval.Kind() == kind {
			return true
		}
	}
	return false
}

func arityErr(expected int, got int, msg string) error {
	if msg == "" {
		return fmt.Errorf("expected %d arguments, got %d", expected, got)
	}

	return fmt.Errorf("expected %d arguments, got %d: %s", expected, got, msg)
}
