package sabre

import (
	"fmt"
	"strings"

	"github.com/spy16/sabre/runtime"
)

var (
	_ runtime.Value     = (*MultiFn)(nil)
	_ runtime.Invokable = (*MultiFn)(nil)
)

// MultiFn represents a multi-arity function definition.
type MultiFn struct {
	Name      string
	Functions []runtime.Fn
}

// Eval returns the multiFn definition itself.
func (multiFn MultiFn) Eval(_ runtime.Runtime) (runtime.Value, error) {
	return multiFn, nil
}

func (multiFn MultiFn) String() string {
	var sb strings.Builder
	sb.WriteString("(defn " + multiFn.Name)
	for _, fn := range multiFn.Functions {
		sb.WriteString("\n  " + fn.String())
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

func (multiFn MultiFn) selectMethod(args []runtime.Value) (runtime.Fn, error) {
	for _, fn := range multiFn.Functions {
		if matchArity(fn, args) {
			return fn, nil
		}
	}
	return runtime.Fn{}, fmt.Errorf("wrong number of args (%d) to '%s'",
		len(args), multiFn.Name)
}

func matchArity(fn runtime.Fn, args []runtime.Value) bool {
	argc := len(args)
	if fn.Variadic {
		return argc >= len(fn.Args)-1
	}
	return argc == len(fn.Args)
}
