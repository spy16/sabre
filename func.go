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

// MacroFn represents a lisp macro.
type MacroFn func(scope Scope, args []Value) (Value, error)

// Eval simply returns the value itself.
func (macro MacroFn) Eval(_ Scope) (Value, error) {
	return macro, nil
}

func (macro MacroFn) String() string {
	return fmt.Sprintf("Macro{%#v}", macro)
}

// Invoke dispatches the call to the underlying Go function.
func (macro MacroFn) Invoke(scope Scope, args ...Value) (Value, error) {
	expanded, err := macro(scope, args)
	if err != nil {
		return nil, err
	}

	return expanded.Eval(scope)
}

// Expand expands the macro and returns the result of expansion.
func (macro MacroFn) Expand(scope Scope, args ...Value) (Value, error) {
	return macro(scope, args)
}
