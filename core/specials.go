package core

import (
	"fmt"
	"reflect"

	"github.com/spy16/sabre/runtime"
)

// Def implements the (def symbol <expr>) special form.
func Def(rt runtime.Runtime, args ...runtime.Value) (runtime.Value, error) {
	if err := VerifyArgCount([]int{2}, len(args)); err != nil {
		return nil, err
	}

	sym, isSym := args[0].(runtime.Symbol)
	if !isSym {
		return nil, fmt.Errorf("first argument to def must be symbol, not '%s'",
			reflect.TypeOf(args[0]))
	}

	val, err := rt.Eval(args[1])
	if err != nil {
		return nil, err
	}

	root := rootRT(rt)
	if err := root.Bind(sym.Value, val); err != nil {
		return nil, err
	}

	return sym, nil
}

func rootRT(rt runtime.Runtime) runtime.Runtime {
	for rt.Parent() != nil {
		rt = rt.Parent()
	}
	return rt
}
