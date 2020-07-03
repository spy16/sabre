package runtime

import (
	"fmt"
	"reflect"
)

var (
	_ Value = (*LinkedList)(nil)
	_ Seq   = (*LinkedList)(nil)
)

// NewSeq returns a new sequence containing given values.
func NewSeq(items ...Value) Seq {
	if len(items) == 0 {
		return Seq((*LinkedList)(nil))
	}
	lst := Seq(&LinkedList{})
	for i := len(items) - 1; i >= 0; i-- {
		lst = Cons(items[i], lst)
	}
	return lst
}

// LinkedList implements Seq using an immutable linked-list.
type LinkedList struct {
	Position
	first Value
	rest  Seq
	count int
}

// Eval evaluates the first item in the list and invokes the resultant first with
// rest of the list as arguments.
func (sl *LinkedList) Eval(rt Runtime) (Value, error) {
	if sl.Count() == 0 {
		return sl, nil
	}

	v, err := rt.Eval(sl.First())
	if err != nil {
		return nil, err
	}

	target, ok := v.(Invokable)
	if !ok {
		return nil, fmt.Errorf("value of type '%s' is not invokable", reflect.TypeOf(v))
	}

	return target.Invoke(rt, toSlice(sl.rest)...)
}

func (sl *LinkedList) String() string { return SeqString(sl, "(", ")", " ") }

// Conj returns a new list with all the items added at the head of the list.
func (sl *LinkedList) Conj(items ...Value) Seq {
	var res Seq
	if sl == nil {
		res = &LinkedList{}
	} else {
		res = sl
	}

	for _, item := range items {
		res = Cons(item, res)
	}
	return res
}

// First returns the head or first item of the list.
func (sl *LinkedList) First() Value {
	if sl == nil {
		return nil
	}
	return sl.first
}

// Next returns the tail of the list.
func (sl *LinkedList) Next() Seq {
	if sl == nil {
		return nil
	}
	return sl.rest
}

// Count returns the number of the list.
func (sl *LinkedList) Count() int {
	if sl == nil {
		return 0
	}
	return sl.count
}
