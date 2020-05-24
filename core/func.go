package core

import (
	"fmt"
	"reflect"
	"strings"
)

// Expand expands the macro invocation form.
func Expand(env Env, form Value) (Value, bool, error) {
	list, ok := form.(*List)
	if !ok || list.Size() == 0 {
		return form, false, nil
	}

	symbol, ok := list.First().(Symbol)
	if !ok {
		return form, false, nil
	}

	target, err := env.Eval(symbol)
	if err != nil || !isMacro(target) {
		return form, false, nil
	}

	mfn := target.(MultiFn)
	v, err := mfn.Expand(env, list.Values[1:])
	return v, true, err
}

// MultiFn represents a multi-arity function or macro definition.
type MultiFn struct {
	Name    string
	IsMacro bool
	Methods []Fn
}

// Source returns source representation of the multi-fn definition.
func (multiFn MultiFn) Source() string {
	var sb strings.Builder
	for _, fn := range multiFn.Methods {
		sb.WriteString("[" + strings.Trim(fn.Source(), "()") + "] ")
	}

	s := multiFn.Name + " " + strings.TrimSpace(sb.String())
	return "(" + strings.TrimSpace(s) + ")"
}

// Invoke dispatches the call to a method based on number of arguments.
func (multiFn MultiFn) Invoke(env Env, args ...Value) (Value, error) {
	if multiFn.IsMacro {
		form, err := multiFn.Expand(env, args)
		if err != nil {
			return nil, err
		}

		return env.Eval(form)
	}

	fn, err := multiFn.selectMethod(args)
	if err != nil {
		return nil, err
	}

	argVals, err := EvalAll(env, args)
	if err != nil {
		return nil, err
	}

	return fn.Invoke(env, argVals...)
}

// Expand executes the macro body and returns the result of the expansion.
func (multiFn MultiFn) Expand(env Env, args []Value) (Value, error) {
	fn, err := multiFn.selectMethod(args)
	if err != nil {
		return nil, err
	}

	if !multiFn.IsMacro {
		return &fn, nil
	}

	return fn.Invoke(env, args...)
}

// Compare returns true if 'v' is also a MultiFn and all methods are
// equivalent.
func (multiFn MultiFn) Compare(v Value) bool {
	other, ok := v.(MultiFn)
	if !ok {
		return false
	}

	sameHeader := (multiFn.Name == other.Name) &&
		(multiFn.IsMacro == other.IsMacro) &&
		(len(multiFn.Methods) == len(other.Methods))
	if !sameHeader {
		return false
	}

	for i, fn1 := range multiFn.Methods {
		fn2 := other.Methods[i]
		if !fn1.Compare(&fn2) {
			return false
		}
	}

	return true
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
	variadicAt := -1
	variadicArity := 0

	for idx, method := range multiFn.Methods {
		if method.Variadic {
			if variadicAt >= 0 {
				return fmt.Errorf("can't have multiple variadic overloads")
			}
			variadicAt = idx
			variadicArity = len(method.Args)
		}
	}

	fixedArities := map[int]struct{}{}
	for idx, method := range multiFn.Methods {
		if method.Variadic {
			continue
		}

		arity := method.minArity()
		if variadicAt >= 0 && idx != variadicAt && arity >= variadicArity {
			return fmt.Errorf("can't have fixed arity overload with more params than variadic")
		}

		if _, exists := fixedArities[arity]; exists {
			return fmt.Errorf("ambiguous arities defined for '%s'", multiFn.Name)
		}
		fixedArities[arity] = struct{}{}
	}

	return nil
}

// Fn represents a function or macro definition.
type Fn struct {
	Args     []string
	Variadic bool
	Body     Value
	Func     func(env Env, args []Value) (Value, error)
}

// Source returns the literal representation of the function.
func (fn Fn) Source() string {
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
func (fn *Fn) Invoke(env Env, args ...Value) (Value, error) {
	if fn.Func != nil {
		return fn.Func(env, args)
	}

	fnEnv := New(env)

	for idx := range fn.Args {
		var argVal Value
		if idx == len(fn.Args)-1 && fn.Variadic {
			argVal = &List{
				Values: args[idx:],
			}
		} else {
			argVal = args[idx]
		}

		_ = fnEnv.Bind(fn.Args[idx], argVal)
	}

	if fn.Body == nil {
		return Nil{}, nil
	}

	return env.Eval(fn.Body)
}

// Compare returns true if 'other' is also a function and has the same
// signature and body.
func (fn *Fn) Compare(v Value) bool {
	other, ok := v.(*Fn)
	if !ok || other == nil {
		return false
	}

	if !reflect.DeepEqual(fn.Args, other.Args) {
		return false
	}

	bothVariadic := (fn.Variadic == other.Variadic)
	noFunc := (fn.Func == nil && other.Func == nil)

	return bothVariadic && noFunc && Compare(fn.Body, other.Body)
}

func (fn Fn) minArity() int {
	if len(fn.Args) > 0 && fn.Variadic {
		return len(fn.Args) - 1
	}
	return len(fn.Args)
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

func isMacro(target Value) bool {
	multiFn, ok := target.(MultiFn)
	return ok && multiFn.IsMacro
}
