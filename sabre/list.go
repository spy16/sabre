package sabre

var (
	_ Value = (*LinkedList)(nil)
	_ Seq   = (*LinkedList)(nil)
)

// Cons returns a new seq with `v` added as the first and `seq` as the rest. Seq
// can be nil as well.
func Cons(v Value, seq Seq) Seq {
	newSeq := &LinkedList{
		first: v,
		rest:  seq,
		count: 1,
	}

	if seq != nil {
		newSeq.count = seq.Count() + 1
	}

	return newSeq
}

// LinkedList implements an immutable Seq using linked-list data structure.
type LinkedList struct {
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
