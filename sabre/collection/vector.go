package collection

import (
	"fmt"
	"reflect"

	"github.com/spy16/sabre/sabre/core"
)

var (
	_ core.Value      = (*Vector)(nil)
	_ core.Seq        = (*Vector)(nil)
	_ core.Comparable = (*Vector)(nil)
)

// Vector represents a list of values. It uses a slice to store the list items.
// Evaluating a Vector evaluates each entry in the Vector and the results are
// returned as another Vector.
type Vector struct {
	core.Position
	items []core.Value
}

// Eval evaluates all the items in the vector and returns a new vector with the
// results of evaluation.
func (vec *Vector) Eval(env core.Env) (core.Value, error) {
	vals, err := core.EvalAll(env, vec.items)
	if err != nil {
		return nil, err
	}
	return &Vector{items: vals}, nil
}

func (vec *Vector) String() string { return core.SeqString(vec, "[", "]", " ") }

// First returns the first item in the vector. Returns nil if the vector is empty.
func (vec *Vector) First() core.Value {
	if len(vec.items) == 0 {
		return nil
	}
	return vec.items[0]
}

// Next returns a sequence containing all the items of the vector except the first
// one. Returns nil if the vector is empty.
func (vec *Vector) Next() core.Seq {
	if len(vec.items) == 0 {
		return nil
	}
	return &Vector{items: append([]core.Value(nil), vec.items[1:]...)}
}

// Cons returns a new vector with 'v' as first item and the current vector as the
// rest.
func (vec *Vector) Cons(v core.Value) core.Seq {
	return &Vector{items: append([]core.Value{v}, vec.items...)}
}

// Conj returns a new vector created by appending 'vals' to this vector.
func (vec *Vector) Conj(vals ...core.Value) core.Seq {
	return &Vector{items: append(vec.items, vals...)}
}

// Count returns the number of items in this vector.
func (vec *Vector) Count() int { return len(vec.items) }

// Compare checks if 'other' is also a vector and then compares each item in the
// vector. Returns true if all compare true.
func (vec *Vector) Compare(other core.Value) bool {
	otherList, ok := other.(*Vector)
	if !ok || vec.Count() != otherList.Count() {
		return false
	}

	for i := 0; i < vec.Count(); i++ {
		if !core.Compare(vec.items[i], otherList.items[i]) {
			return false
		}
	}

	return true
}

func (vec *Vector) toIndex(key core.Value) (int, error) {
	idx, ok := key.(core.Int64)
	if !ok {
		return 0, fmt.Errorf("key must be integer, not '%s'", reflect.TypeOf(key))
	} else if idx < 0 || int(idx) >= vec.Count() {
		return 0, fmt.Errorf("index out of bounds")
	}

	return int(idx), nil
}

// VectorReader implements the reader macro for reading vector from source.
func VectorReader(rd *core.Reader, _ rune) (core.Value, error) {
	const vecEnd = ']'

	pi := rd.Position()
	forms, err := rd.Container(vecEnd, "vector")
	if err != nil {
		return nil, err
	}

	return &Vector{
		items:    forms,
		Position: pi,
	}, nil
}
