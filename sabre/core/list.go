package core

import (
	"fmt"
	"reflect"
)

var (
	_ Value      = (*List)(nil)
	_ Seq        = (*List)(nil)
	_ Comparable = (*List)(nil)
)

// List represents a list of values. List can be backed by any Seq implementation.
// Evaluating a list leads to invocation of result of evaluation of first entry in
// the list.
type List struct {
	Position
	Items []Value
}

// Eval evaluates the first item in the list and invokes the resultant value with
// rest of the list as arguments.
func (sl *List) Eval(env Env) (Value, error) {
	if sl.Count() == 0 {
		return sl, nil
	}

	v, err := env.Eval(sl.First())
	if err != nil {
		return nil, err
	}

	target, ok := v.(Invokable)
	if !ok {
		return nil, fmt.Errorf("value of type '%s' is not invokable", reflect.TypeOf(target))
	}

	var args []Value
	ForEach(sl, func(item Value) bool {
		args = append(args, item)
		return false
	})

	return target.Invoke(env, args...)
}

func (sl List) String() string {
	return SeqString(&sl, "(", ")", " ")
}

// Compare checks if the other value is a list and if it is, compares each item
// in the list. Returns true if all match.
func (sl *List) Compare(other Value) bool {
	otherList, ok := other.(*List)
	if !ok || sl.Count() != otherList.Count() {
		return false
	}

	var prev Value
	isSame := true
	ForEach(sl, func(item Value) bool {
		if prev == nil {
			return false
		}
		isSame = isSame && Compare(prev, item)
		return !isSame
	})

	return isSame
}

// Count returns the number of items in the list.
func (sl *List) Count() int { return len(sl.Items) }

// Cons returns a new list with the given item added to the front.
func (sl *List) Cons(v Value) Seq {
	return &List{Items: append([]Value{v}, sl.Items...)}
}

// Conj returns a new list with the given vals appended.
func (sl *List) Conj(vals ...Value) Seq {
	return &List{Items: append(sl.Items, vals...)}
}

// First returns the first item in the list or nil if list is empty.
func (sl *List) First() Value {
	if len(sl.Items) == 0 {
		return nil
	}
	return sl.Items[0]
}

// Next returns a list containing all but first item in the list. Returns
// nil if the list is empty.
func (sl *List) Next() Seq {
	if len(sl.Items) == 0 {
		return nil
	}
	return &List{Items: append([]Value(nil), sl.Items[1:]...)}
}
