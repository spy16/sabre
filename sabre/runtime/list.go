package runtime

import (
	"fmt"
	"reflect"
)

var (
	_ Value = (*linkedList)(nil)
	_ Seq   = (*linkedList)(nil)
)

// NewSeq returns a new sequence containing given values.
func NewSeq(items ...Value) Seq {
	if len(items) == 0 {
		return Seq((*linkedList)(nil))
	}
	lst := Seq(&linkedList{})
	for i := len(items) - 1; i >= 0; i-- {
		lst = lst.Cons(items[i])
	}
	return lst
}

// linkedList implements Seq using an immutable linked-list.
type linkedList struct {
	Position
	value Value
	next  *linkedList
	count int
}

// Eval evaluates the first item in the list and invokes the resultant value with
// rest of the list as arguments.
func (sl *linkedList) Eval(rt Runtime) (Value, error) {
	if sl.Count() == 0 {
		return sl, nil
	}

	v, err := rt.Eval(sl.First())
	if err != nil {
		return nil, err
	}

	target, ok := v.(Invokable)
	if !ok {
		return nil, fmt.Errorf("value of type '%s' is not invokable", reflect.TypeOf(target))
	}

	var args []Value
	ForEach(sl.next, func(item Value) bool {
		args = append(args, item)
		return false
	})

	return target.Invoke(rt, args...)
}

func (sl *linkedList) String() string { return SeqString(sl, "(", ")", " ") }

// Cons returns a new list with 'v' added as head and current list as tail.
func (sl *linkedList) Cons(v Value) Seq {
	newSeq := &linkedList{
		value: v,
		next:  sl,
		count: 1,
	}

	if sl != nil {
		newSeq.count = sl.count + 1
		newSeq.Position = sl.Position
	}

	return newSeq
}

// Conj returns a new list with all the items added at the head of the list.
func (sl *linkedList) Conj(items ...Value) Seq {
	if sl == nil {
		sl = &linkedList{}
	}

	res := Seq(sl)
	for _, item := range items {
		res = res.Cons(item)
	}
	return res
}

// First returns the head or first item of the list.
func (sl *linkedList) First() Value {
	if sl == nil {
		return nil
	}
	return sl.value
}

// Next returns the tail of the list.
func (sl *linkedList) Next() Seq {
	if sl == nil {
		return nil
	}
	return sl.next
}

// Count returns the number of the list.
func (sl *linkedList) Count() int {
	if sl == nil {
		return 0
	}
	return sl.count
}
