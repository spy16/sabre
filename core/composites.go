package core

import (
	"fmt"
	"reflect"
	"strings"
)

// List represents a list of forms/vals. Evaluating a list leads to a function
// invocation.
type List struct {
	Values
	Position
}

// Eval performs an invocation. Result of evaluating the first item in the list
// must be Invokable. Result of Eval will be the result of invoke.
func (lf *List) Eval(env Env) (Value, error) {
	if lf.Size() == 0 {
		return lf, nil
	}

	target, err := env.Eval(lf.Values[0])
	if err != nil {
		return nil, err
	}

	invokable, ok := target.(Invokable)
	if !ok {
		return nil, fmt.Errorf(
			"cannot invoke value of type '%s'", reflect.TypeOf(target),
		)
	}

	return invokable.Invoke(env, lf.Values[1:]...)
}

// Source returns the source representation of the list.
func (lf List) Source() string {
	return containerString(lf.Values, "(", ")", " ")
}

// Vector represents a list of values. Unlike List type, evaluation of vector
// does not lead to function invoke.
type Vector struct {
	Values
	Position
}

// Eval evaluates each value in the vector form and returns resultant values
// as new vector.
func (vf Vector) Eval(env Env) (Value, error) {
	vals, err := EvalAll(env, vf.Values)
	if err != nil {
		return nil, err
	}
	return Vector{Values: vals}, nil
}

// Invoke of a vector performs a index lookup. Only arity 1 is allowed and
// should be an integer value to be used as index.
func (vf Vector) Invoke(env Env, args ...Value) (Value, error) {
	vals, err := EvalAll(env, args)
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

// Source returns the source representation of the vector by iteratively
// finding source for all contained forms.
func (vf Vector) Source() string {
	return containerString(vf.Values, "[", "]", " ")
}

// Set represents a list of unique values. (Experimental)
type Set struct {
	Values
	Position
}

// Eval evaluates each value in the set form and returns the resultant
// values as new set.
func (set Set) Eval(env Env) (Value, error) {
	vals, err := EvalAll(env, set.Uniq())
	if err != nil {
		return nil, err
	}
	return Set{Values: Values(vals).Uniq()}, nil
}

// Source returns the source representation for the set.
func (set Set) Source() string {
	return containerString(set.Values, "#{", "}", " ")
}

// Valid checks if the set does not contain any duplicates.
// TODO: Remove this naive solution
func (set Set) Valid() bool {
	s := map[string]struct{}{}

	for _, v := range set.Values {
		str := v.Source()
		if _, found := s[str]; found {
			return false
		}
		s[str] = struct{}{}
	}

	return true
}

// HashMap represents a container for key-value pairs.
type HashMap struct {
	Position
	Data map[Value]Value
}

// Eval evaluates all keys and values and returns a new HashMap containing
// the evaluated values.
func (hm *HashMap) Eval(env Env) (Value, error) {
	res := &HashMap{Data: map[Value]Value{}}
	for k, v := range hm.Data {
		key, err := env.Eval(k)
		if err != nil {
			return nil, err
		}

		val, err := env.Eval(v)
		if err != nil {
			return nil, err
		}

		res.Data[key] = val
	}

	return res, nil
}

// Source returns the source representation of the hashmap.
func (hm *HashMap) Source() string {
	var fields []Value
	for k, v := range hm.Data {
		fields = append(fields, k, v)
	}
	return containerString(fields, "{", "}", " ")
}

// Get returns the value associated with the given key if found.
// Returns def otherwise.
func (hm *HashMap) Get(key Value, def Value) Value {
	if !IsHashable(key) {
		return def
	}

	v, found := hm.Data[key]
	if !found {
		return def
	}

	return v
}

// Set sets/updates the value associated with the given key.
func (hm *HashMap) Set(key, val Value) error {
	if !IsHashable(key) {
		return fmt.Errorf("value of type '%s' is not hashable", key)
	}

	if hm.Data == nil {
		hm.Data = map[Value]Value{}
	}

	hm.Data[key] = val
	return nil
}

// Keys returns all the keys in the hashmap.
func (hm *HashMap) Keys() Values {
	var res []Value
	for k := range hm.Data {
		res = append(res, k)
	}
	return res
}

// Values returns all the values in the hashmap.
func (hm *HashMap) Values() Values {
	var res []Value
	for _, v := range hm.Data {
		res = append(res, v)
	}
	return res
}

// Module represents a group of forms. Evaluating a module leads to evaluation
// of each form in order and result will be the result of last evaluation.
type Module []Value

// Eval evaluates all the vals in the module body and returns the result of the
// last evaluation.
func (mod Module) Eval(env Env) (Value, error) {
	if len(mod) == 0 {
		return Nil{}, nil
	}

	res, err := EvalAll(env, mod)
	if err != nil {
		return nil, err
	}
	return res[len(res)-1], nil
}

// Compare returns true if the 'v' is also a module and all forms in the
// module are equivalent.
func (mod Module) Compare(v Value) bool {
	otherMod, ok := v.(Module)
	if !ok {
		return false
	}

	if len(mod) != len(otherMod) {
		return false
	}

	for i := range mod {
		if !Compare(mod[i], otherMod[i]) {
			return false
		}
	}

	return true
}

// Source iteratively converts the entire module into source representation.
func (mod Module) Source() string {
	return containerString(mod, "", "\n", "\n")
}

// IsHashable returns true if the value can be used as key for a hashmap.
func IsHashable(v Value) bool {
	switch v.(type) {
	case String, Int64, Float64, Nil, Character, Keyword:
		return true

	default:
		return false
	}
}

func containerString(vals []Value, begin, end, sep string) string {
	parts := make([]string, len(vals))
	for i, expr := range vals {
		parts[i] = fmt.Sprintf("%v", expr)
	}
	return begin + strings.Join(parts, sep) + end
}
