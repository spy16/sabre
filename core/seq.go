package core

// Seq implementations represent a sequence/list of values.
type Seq interface {
	Value

	// First should return first value of the sequence or nil if the
	// sequence is empty.
	First() Value
	// Next should return the remaining sequence when the first value
	// is excluded.
	Next() Seq
	// Cons should add the value to the beginning of the sequence and
	// return the new sequence.
	Cons(v Value) Seq
	// Conj should join the given values to the sequence and return a
	// new sequence.
	Conj(vals ...Value) Seq
}

// Values represents a list of values and implements the Seq interface.
type Values []Value

// First returns the first value in the list if the list is not empty.
// Returns Nil{} otherwise.
func (vals Values) First() Value {
	if len(vals) == 0 {
		return nil
	}
	return vals[0]
}

// Next returns a new sequence containing values after the first one. If
// there are no values to create a next sequence, returns nil.
func (vals Values) Next() Seq {
	if len(vals) <= 1 {
		return nil
	}
	return Values(vals[1:])
}

// Cons returns a new sequence where 'v' is prepended to the values.
func (vals Values) Cons(v Value) Seq {
	return append(Values{v}, vals...)
}

// Conj returns a new sequence where 'v' is appended to the values.
func (vals Values) Conj(args ...Value) Seq {
	return append(vals, args...)
}

// Size returns the number of items in the list.
func (vals Values) Size() int { return len(vals) }

// Source returns list representation of the value list.
func (vals Values) Source() string {
	return containerString(vals, "(", ")", " ")
}

// Compare compares the values in this sequence to the other sequence.
// other sequence will be realized for comparison.
func (vals Values) Compare(v Value) bool {
	other, ok := v.(Seq)
	if !ok {
		return false
	}

	if s, hasSize := other.(interface {
		Size() int
	}); hasSize {
		if vals.Size() != s.Size() {
			return false
		}
	}

	var this Seq = vals
	isEqual := true
	for this != nil && other != nil {
		v1, v2 := this.First(), other.First()
		isEqual = isEqual && Compare(v1, v2)
		if !isEqual {
			break
		}

		this = this.Next()
		other = other.Next()
	}

	return isEqual && (this == nil && other == nil)
}

// Uniq removes all the duplicates from the given value array.
// TODO: remove this naive implementation
func (vals Values) Uniq() []Value {
	var result []Value

	hashSet := map[string]struct{}{}
	for _, v := range vals {
		src := v.Source()
		if _, found := hashSet[src]; !found {
			hashSet[src] = struct{}{}
			result = append(result, v)
		}
	}

	return result
}
