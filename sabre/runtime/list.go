package runtime

import (
	"fmt"
	"reflect"
)

var (
	_ Value = (*SliceList)(nil)
	_ Seq   = (*SliceList)(nil)
)

// SliceList represents a list of values. SliceList is backed by a native Go slice. Eval
// leads to invocation of result of evaluation of first entry in the list.
type SliceList struct {
	Position
	Items []Value
}

// Eval evaluates the first item in the list and invokes the resultant value with
// rest of the list as arguments.
func (sl *SliceList) Eval(env Runtime) (Value, error) {
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

func (sl SliceList) String() string { return SeqString(&sl, "(", ")", " ") }

// Count returns the number of items in the list.
func (sl *SliceList) Count() int { return len(sl.Items) }

// Cons returns a new list with the given item added to the front.
func (sl *SliceList) Cons(v Value) Seq {
	return &SliceList{Items: append([]Value{v}, sl.Items...)}
}

// Conj returns a new list with the given vals appended.
func (sl *SliceList) Conj(vals ...Value) Seq {
	return &SliceList{Items: append(sl.Items, vals...)}
}

// First returns the first item in the list or nil if list is empty.
func (sl *SliceList) First() Value {
	if len(sl.Items) == 0 {
		return nil
	}
	return sl.Items[0]
}

// Next returns a list containing all but first item in the list. Returns
// nil if the list is empty.
func (sl *SliceList) Next() Seq {
	if len(sl.Items) == 0 {
		return nil
	}
	return &SliceList{Items: append([]Value(nil), sl.Items[1:]...)}
}
