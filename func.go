package sabre

import "fmt"

// GoFunc implements Invokable using a Go function value.
type GoFunc func(scope Scope, args []Value) (Value, error)

// Eval simply returns the value itself.
func (goFn GoFunc) Eval(_ Scope) (Value, error) {
	return goFn, nil
}

func (goFn GoFunc) String() string {
	return fmt.Sprintf("GoFunc{%#v}", goFn)
}

// Invoke dispatches the call to the underlying Go function.
func (goFn GoFunc) Invoke(scope Scope, args ...Value) (Value, error) {
	return goFn(scope, args)
}
