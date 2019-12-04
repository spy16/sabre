package sabre

import "fmt"

// Vector represents a list of values. Unlike List type, evaluation of
// vector does not lead to function invoke.
type Vector []Value

// Eval evaluates each value in the vector form and returns the resultant
// values as new vector.
func (vf Vector) Eval(scope Scope) (Value, error) {
	vals, err := evalValueList(scope, vf)
	if err != nil {
		return nil, err
	}

	return Vector(vals), nil
}

// Invoke of a vector performs a index lookup. Only arity 1 is allowed
// and should be an integer value to be used as index.
func (vf Vector) Invoke(args []Value) (Value, error) {
	if len(args) != 1 {
		return nil, arityErr(1, len(args), "")
	}

	key := newValue(args[0])
	index, err := key.ToInt64()
	if err != nil {
		return nil, err
	}

	if int(index) >= len(vf) {
		return nil, fmt.Errorf("index out of bounds")
	}

	return vf[index], nil
}

func (vf Vector) String() string {
	return containerString(vf, "[", "]", " ")
}

func readVector(rd *Reader, _ rune) (Value, error) {
	forms, err := readContainer(rd, '[', ']', "vector")
	if err != nil {
		return nil, err
	}

	return Vector(forms), nil
}

func arityErr(expected int, got int, msg string) error {
	if msg == "" {
		return fmt.Errorf("expected %d arguments, got %d", expected, got)
	}

	return fmt.Errorf("expected %d arguments, got %d: %s", expected, got, msg)
}
