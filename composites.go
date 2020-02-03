package sabre

import (
	"fmt"
	"reflect"
	"strings"
)

// List represents an list of forms/vals. Evaluating a list leads to a
// function invocation.
type List struct {
	Values

	// special is set in case if the list represents invocation of
	// a special form such as def, fn* etc.
	special func(scope Scope) (Value, error)
}

// Eval performs an invocation.
func (lf *List) Eval(scope Scope) (Value, error) {
	if lf.Size() == 0 {
		return &List{}, nil
	}

	if lf.special != nil {
		return lf.special(scope)
	}

	target, err := lf.Values[0].Eval(scope)
	if err != nil {
		return nil, err
	}

	fn, ok := target.(Invokable)
	if !ok {
		return nil, fmt.Errorf("cannot invoke value of type '%s'", reflect.TypeOf(target))
	}

	return fn.Invoke(scope, lf.Values[1:]...)
}

func (lf List) String() string {
	return containerString(lf.Values, "(", ")", " ")
}

func (lf *List) parseSpecial(scope Scope) error {
	if lf.Size() == 0 {
		return nil
	}

	sym, isSymbol := lf.Values[0].(Symbol)
	if !isSymbol {
		return nil
	}

	special, isSpecial := specialForms[sym.Value]
	if !isSpecial {
		return nil
	}

	expr, err := special(scope, lf.Values[1:])
	if err != nil {
		return err
	}

	lf.special = expr
	return nil
}

// Vector represents a list of values. Unlike List type, evaluation of
// vector does not lead to function invoke.
type Vector struct {
	Values
}

// Eval evaluates each value in the vector form and returns the resultant
// values as new vector.
func (vf Vector) Eval(scope Scope) (Value, error) {
	vals, err := evalValueList(scope, vf.Values)
	if err != nil {
		return nil, err
	}

	return Vector{Values: vals}, nil
}

// Invoke of a vector performs a index lookup. Only arity 1 is allowed
// and should be an integer value to be used as index.
func (vf Vector) Invoke(scope Scope, args ...Value) (Value, error) {
	vals, err := evalValueList(scope, args)
	if err != nil {
		return nil, err
	}

	if len(vals) != 1 {
		return nil, fmt.Errorf("call requires exactly 1 argument, got %d", len(vals))
	}

	index, isInt := vals[0].(Int64)
	if !isInt {
		return nil, fmt.Errorf("key must be integer")
	}

	if int(index) >= len(vf.Values) {
		return nil, fmt.Errorf("index out of bounds")
	}

	return vf.Values[index], nil
}

func (vf Vector) String() string {
	return containerString(vf.Values, "[", "]", " ")
}

// Set represents a list of unique values. (Experimental)
type Set struct {
	Values
}

// Eval evaluates each value in the set form and returns the resultant
// values as new set.
func (set Set) Eval(scope Scope) (Value, error) {
	vals, err := evalValueList(scope, set.Values)
	if err != nil {
		return nil, err
	}

	return Set{Values: uniq(vals)}, nil
}

func (set Set) String() string {
	return containerString(set.Values, "#{", "}", " ")
}

func (set Set) valid() bool {
	// TODO: Remove this naive solution
	s := map[string]struct{}{}

	for _, v := range set.Values {
		str := v.String()
		if _, found := s[str]; found {
			return false
		}
		s[v.String()] = struct{}{}
	}

	return true
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
		return Nil{}, nil
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

func containerString(vals []Value, begin, end, sep string) string {
	parts := make([]string, len(vals))
	for i, expr := range vals {
		parts[i] = fmt.Sprintf("%v", expr)
	}
	return begin + strings.Join(parts, sep) + end
}

func uniq(items []Value) []Value {
	// TODO: remove this naive implementation
	vs := map[string]Value{}
	for _, v := range items {
		s := v.String()
		vs[s] = v
	}

	var set []Value
	for _, v := range vs {
		set = append(set, v)
	}

	return set
}
