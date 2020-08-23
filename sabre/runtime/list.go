package runtime

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

// LinkedList implements an immutable Seq using linked-list data structure.
type LinkedList struct {
	Position
	first Value
	rest  Seq
	count int
}

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

func (sl *LinkedList) String() string { return SeqString(sl, "(", ")", " ") }
