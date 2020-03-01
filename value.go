package sabre

import "reflect"

// Value represents data/forms in sabre. This includes those emitted by
// Reader, values obtained as result of an evaluation etc.
type Value interface {
	// String should return the LISP representation of the value.
	String() string
	// Eval should evaluate this value against the scope and return
	// the resultant value or an evaluation error.
	Eval(scope Scope) (Value, error)
}

// Invokable represents any value that supports invocation. Vector, Fn
// etc support invocation.
type Invokable interface {
	Value
	Invoke(scope Scope, args ...Value) (Value, error)
}

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

// Compare compares two values in an identity independent manner. If
// v1 has `Compare(Value) bool` method, the comparison is delegated to
// it as `v1.Compare(v2)`.
func Compare(v1, v2 Value) bool {
	if (v1 == nil && v2 == nil) ||
		(v1 == (Nil{}) && v2 == (Nil{})) {
		return true
	}

	if cmp, ok := v1.(comparable); ok {
		return cmp.Compare(v2)
	}

	return reflect.DeepEqual(v1, v2)
}

// comparable can be implemented by Value types to support comparison.
// See Compare().
type comparable interface {
	Value
	Compare(other Value) bool
}

// Values represents a list of values and implements the Seq interface.
type Values []Value

// Eval returns itself.
func (vals Values) Eval(_ Scope) (Value, error) { return vals, nil }

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
	return &List{Values: Values(vals[1:])}
}

// Cons returns a new sequence where 'v' is prepended to the values.
func (vals Values) Cons(v Value) Seq {
	return &List{Values: append(Values{v}, vals...)}
}

// Conj returns a new sequence where 'v' is appended to the values.
func (vals Values) Conj(args ...Value) Seq {
	return &List{Values: append(vals, args...)}
}

// Size returns the number of items in the list.
func (vals Values) Size() int { return len(vals) }

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
		src := v.String()
		if _, found := hashSet[src]; !found {
			hashSet[src] = struct{}{}
			result = append(result, v)
		}
	}

	return result
}

func (vals Values) String() string {
	return containerString(vals, "(", ")", " ")
}

func evalValueList(scope Scope, vals []Value) ([]Value, error) {
	var result []Value

	for _, arg := range vals {
		v, err := arg.Eval(scope)
		if err != nil {
			return nil, newEvalErr(arg, err)
		}

		result = append(result, v)
	}

	return result, nil
}
