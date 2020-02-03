package sabre

import "fmt"

// MultiFn represents a multi-arity function or macro definition.
type MultiFn struct {
	Name    string
	IsMacro bool
	Methods []Fn
}

// Eval returns the multiFn definition itself.
func (multiFn MultiFn) Eval(_ Scope) (Value, error) {
	return multiFn, nil
}

func (multiFn MultiFn) String() string {
	return fmt.Sprintf("MultiFn{name=%s}", multiFn.Name)
}

// Invoke dispatches the call to a method based on number of arguments.
func (multiFn MultiFn) Invoke(scope Scope, args ...Value) (Value, error) {
	fn, err := multiFn.selectMethod(args)
	if err != nil {
		return nil, err
	}

	if multiFn.IsMacro {
		v, err := fn.Invoke(scope, args)
		if err != nil {
			return nil, err
		}

		return v.Eval(scope)
	}

	argVals, err := evalValueList(scope, args)
	if err != nil {
		return nil, err
	}

	return fn.Invoke(scope, argVals)
}

func (multiFn MultiFn) selectMethod(args []Value) (Fn, error) {
	for _, fn := range multiFn.Methods {
		if fn.matchArity(args) {
			return fn, nil
		}
	}

	return Fn{}, fmt.Errorf("wrong number of args (%d) to '%s'",
		len(args), multiFn.Name)
}

// Fn represents a function or macro definition.
type Fn struct {
	Args     []string
	Variadic bool
	Body     Value
	Func     Invokable
}

// Invoke executes the function with given arguments.
func (fn Fn) Invoke(scope Scope, args []Value) (Value, error) {
	if fn.Func != nil {
		return fn.Func.Invoke(scope, args...)
	}

	fnScope := NewScope(scope)

	for idx := range fn.Args {
		var argVal Value
		if idx == len(fn.Args)-1 && fn.Variadic {
			argVal = &List{
				Values: args[idx:],
			}
		} else {
			argVal = args[idx]
		}

		_ = fnScope.Bind(fn.Args[idx], argVal)
	}

	if fn.Body == nil {
		return Nil{}, nil
	}

	return fn.Body.Eval(fnScope)
}

func (fn Fn) matchArity(args []Value) bool {
	argc := len(args)

	if fn.Variadic {
		return argc >= len(fn.Args)-1
	}

	return argc == len(fn.Args)
}

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
