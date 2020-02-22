package sabre

import (
	"fmt"
	"reflect"
	"strings"
)

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
	var sb strings.Builder
	for _, fn := range multiFn.Methods {
		sb.WriteString("[" + strings.Trim(fn.String(), "()") + "] ")
	}

	s := "(" + strings.TrimSpace(multiFn.Name+" "+strings.TrimSpace(sb.String())) + ")"
	return s
}

// Invoke dispatches the call to a method based on number of arguments.
func (multiFn MultiFn) Invoke(scope Scope, args ...Value) (Value, error) {
	fn, err := multiFn.selectMethod(args)
	if err != nil {
		return nil, err
	}

	if multiFn.IsMacro {
		v, err := fn.Invoke(scope, args...)
		if err != nil {
			return nil, err
		}

		return Eval(scope, v)
	}

	argVals, err := evalValueList(scope, args)
	if err != nil {
		return nil, err
	}

	return fn.Invoke(scope, argVals...)
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

func (multiFn *MultiFn) validate() error {
	return nil
}

// Fn represents a function or macro definition.
type Fn struct {
	Args     []string
	Variadic bool
	Body     Value
	Func     func(scope Scope, args []Value) (Value, error)
}

// Eval returns the function itself.
func (fn *Fn) Eval(_ Scope) (Value, error) {
	return fn, nil
}

func (fn Fn) String() string {
	var sb strings.Builder

	for i, arg := range fn.Args {
		if i == len(fn.Args)-1 && fn.Variadic {
			sb.WriteString(" & " + arg)
		} else {
			sb.WriteString(arg + " ")
		}
	}

	return "(" + strings.TrimSpace(sb.String()) + ")"
}

// Invoke executes the function with given arguments.
func (fn *Fn) Invoke(scope Scope, args ...Value) (Value, error) {
	if fn.Func != nil {
		return fn.Func(scope, args)
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

	return Eval(fnScope, fn.Body)
}

func (fn Fn) matchArity(args []Value) bool {
	argc := len(args)

	if fn.Variadic {
		return argc >= len(fn.Args)-1
	}

	return argc == len(fn.Args)
}

func (fn *Fn) parseArgSpec(spec Value) error {
	vec, isVector := spec.(Vector)
	if !isVector {
		return fmt.Errorf("argument spec must be a vector of symbols, not '%s'",
			reflect.TypeOf(spec))
	}

	argNames, err := toArgNames(vec.Values)
	if err != nil {
		return err
	}

	fn.Variadic, err = checkVariadic(argNames)
	if err != nil {
		return err
	}

	if fn.Variadic {
		argc := len(argNames)
		fn.Args = append(argNames[:argc-2], argNames[argc-1])
	} else {
		fn.Args = argNames
	}

	return nil
}

func checkVariadic(args []string) (bool, error) {
	for i, arg := range args {
		if arg != "&" {
			continue
		}

		if i > len(args)-2 {
			return false, fmt.Errorf("expecting one more symbol after '&'")
		} else if i < len(args)-2 {
			return false, fmt.Errorf("expecting only one symbol after '&'")
		}

		return true, nil
	}

	return false, nil
}

func toArgNames(vals []Value) ([]string, error) {
	var names []string

	for i, v := range vals {
		sym, isSymbol := v.(Symbol)
		if !isSymbol {
			return nil, fmt.Errorf(
				"expecting symbol at '%d', not '%s'",
				i, reflect.TypeOf(v),
			)
		}

		names = append(names, sym.Value)
	}

	return names, nil
}
