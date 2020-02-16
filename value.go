package sabre

// Value represents data/forms in sabre. This includes those emitted by
// Reader, values obtained as result of an evaluation etc.
type Value interface {
	// Eval should evaluate this value against the scope and return
	// the resultant value or an evaluation error.
	Eval(scope Scope) (Value, error)

	// String should return the LISP representation of the value.
	String() string
}

// Invokable represents any value that supports invocation. Vector, Fn
// etc support invocation.
type Invokable interface {
	Invoke(scope Scope, args ...Value) (Value, error)
}

// Seq implementations represent a sequence/list of values.
type Seq interface {
	Value
	First() Value
	Next() Seq
	Cons(v Value) Seq
	Conj(vals ...Value) Seq
}

// Values represents a list of values and implements the Seq interface.
type Values []Value

// Eval returns itself.
func (vals Values) Eval(_ Scope) (Value, error) {
	return vals, nil
}

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
func (vals Values) Size() int {
	return len(vals)
}

func (vals Values) String() string {
	return containerString(vals, "(", ")", " ")
}
