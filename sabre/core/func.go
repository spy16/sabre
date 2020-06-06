package core

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spy16/sabre/sabre/runtime"
)

var (
	_ runtime.Value     = (*MultiFn)(nil)
	_ runtime.Invokable = (*MultiFn)(nil)

	_ runtime.Value     = (*Fn)(nil)
	_ runtime.Invokable = (*Fn)(nil)
)

// MultiFn represents a multi-arity function definition.
type MultiFn struct {
	Name      string
	Functions []Fn
}

// Fn represents a function definition.
type Fn struct {
	Args     []string
	Variadic bool
	Body     runtime.Value
}

// Eval returns the multiFn definition itself.
func (multiFn MultiFn) Eval(_ runtime.Runtime) (runtime.Value, error) {
	return multiFn, nil
}

func (multiFn MultiFn) String() string {
	var sb strings.Builder
	sb.WriteString("(defn " + multiFn.Name)
	for _, fn := range multiFn.Functions {
		sb.WriteString("\n  " + fn.stringWithPrefix(""))
	}
	sb.WriteRune(')')
	return sb.String()
}

// Invoke dispatches the call to a method based on number of arguments.
func (multiFn MultiFn) Invoke(rt runtime.Runtime, args ...runtime.Value) (runtime.Value, error) {
	fn, err := multiFn.selectMethod(args)
	if err != nil {
		return nil, err
	}

	argVals, err := runtime.EvalAll(rt, args)
	if err != nil {
		return nil, err
	}

	return fn.Invoke(rt, argVals...)
}

// Equals returns true if 'other' is also a MultiFn and all functions are
// equivalent.
func (multiFn MultiFn) Equals(other runtime.Value) bool {
	otherMultiFn, ok := other.(MultiFn)
	if !ok {
		return false
	}

	sameHeader := (multiFn.Name == otherMultiFn.Name) &&
		(len(multiFn.Functions) == len(otherMultiFn.Functions))
	if !sameHeader {
		return false
	}

	for i, fn1 := range multiFn.Functions {
		fn2 := otherMultiFn.Functions[i]
		if !fn1.Equals(&fn2) {
			return false
		}
	}

	return true
}

func (multiFn MultiFn) selectMethod(args []runtime.Value) (Fn, error) {
	for _, fn := range multiFn.Functions {
		if fn.matchArity(args) {
			return fn, nil
		}
	}

	return Fn{}, fmt.Errorf("wrong number of args (%d) to '%s'",
		len(args), multiFn.Name)
}

// Eval returns the function itself.
func (fn *Fn) Eval(_ runtime.Runtime) (runtime.Value, error) { return fn, nil }

func (fn Fn) String() string { return fn.stringWithPrefix("fn") }

// Invoke executes the function with given arguments.
func (fn *Fn) Invoke(rt runtime.Runtime, args ...runtime.Value) (runtime.Value, error) {
	fnEnv := runtime.New(rt)

	for idx := range fn.Args {
		var argVal runtime.Value
		if idx == len(fn.Args)-1 && fn.Variadic {
			argVal = runtime.NewSeq(args[idx:]...)
		} else {
			argVal = args[idx]
		}

		_ = fnEnv.Bind(fn.Args[idx], argVal)
	}

	if fn.Body == nil {
		return runtime.Nil{}, nil
	}

	return fnEnv.Eval(fn.Body)
}

// Equals returns true if 'other' is also an Fn value and has the same
// signature and body.
func (fn *Fn) Equals(other runtime.Value) bool {
	otherFn, ok := other.(*Fn)
	if !ok || otherFn == nil {
		return false
	}

	sameArgs := reflect.DeepEqual(fn.Args, otherFn.Args)
	bothVariadic := (fn.Variadic == otherFn.Variadic)
	return bothVariadic && sameArgs && runtime.Equals(fn.Body, otherFn.Body)
}

func (fn Fn) stringWithPrefix(prefix string) string {
	var sb strings.Builder
	sb.WriteRune('(')

	if prefix != "" {
		sb.WriteString(prefix + " ")
	}

	sb.WriteRune('[')
	if len(fn.Args) > 0 {
		if fn.Variadic {
			sb.WriteString(strings.Join(fn.Args[:len(fn.Args)-1], " "))
			sb.WriteString(" & " + fn.Args[len(fn.Args)-1])
		} else {
			sb.WriteString(strings.Join(fn.Args, " "))
		}
	}
	sb.WriteString("] ")
	if fn.Body != nil {
		sb.WriteString(fn.Body.String())
	}
	sb.WriteRune(')')

	return sb.String()
}

func (fn Fn) matchArity(args []runtime.Value) bool {
	argc := len(args)
	if fn.Variadic {
		return argc >= len(fn.Args)-1
	}
	return argc == len(fn.Args)
}
