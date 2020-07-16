package sabre

import (
	"fmt"
	"reflect"

	"github.com/spy16/sabre/core"
	"github.com/spy16/sabre/runtime"
)

func letForm(rt runtime.Runtime, args ...runtime.Value) (runtime.Value, error) {
	if len(args) == 0 {
		return runtime.Nil{}, nil
	}

	bindings, isVec := args[0].(runtime.Vector)
	if !isVec || bindings.Count()%2 != 0 {
		return nil, fmt.Errorf("first arg to let form must be a vector with even forms")
	}

	letRT := runtime.New(rt)
	for i := 0; i < bindings.Count(); i += 2 {
		v0, err := bindings.EntryAt(i)
		if err != nil {
			return nil, err
		}
		v1, err := bindings.EntryAt(i + 1)
		if err != nil {
			return nil, err
		}

		sym, isSym := v0.(runtime.Symbol)
		if !isSym {
			return nil, fmt.Errorf("form at %d must be a symbol, not '%s'",
				i, reflect.TypeOf(v0))
		}

		val, err := rt.Eval(v1)
		if err != nil {
			return nil, err
		}

		if err := letRT.Bind(sym.Value, val); err != nil {
			return nil, err
		}
	}

	return letRT.Eval(core.Module(args[1:]))
}

func defForm(rt runtime.Runtime, args ...runtime.Value) (runtime.Value, error) {
	if err := core.VerifyArgCount([]int{2}, len(args)); err != nil {
		return nil, err
	}

	sym, isSym := args[0].(runtime.Symbol)
	if !isSym {
		return nil, fmt.Errorf("first argument to def must be Symbol, not '%s'",
			reflect.TypeOf(args[0]))
	}

	v, err := rt.Eval(args[1])
	if err != nil {
		return nil, err
	}

	return sym, rootRT(rt).Bind(sym.Value, v)
}

func doForm(rt runtime.Runtime, args ...runtime.Value) (runtime.Value, error) {
	if len(args) == 0 {
		return runtime.Nil{}, nil
	}

	res, err := runtime.EvalAll(rt, args)
	if err != nil {
		return nil, err
	}
	return res[len(res)-1], nil
}

func condForm(rt runtime.Runtime, args ...runtime.Value) (runtime.Value, error) {
	if len(args)%2 != 0 {
		return nil, fmt.Errorf("cond requires even number of forms")
	}

	for i := 0; i < len(args); i += 2 {
		test, err := rt.Eval(args[i])
		if err != nil {
			return nil, err
		}

		if isTruthy(test) {
			return rt.Eval(args[i+1])
		}
	}

	return runtime.Nil{}, nil
}

func isTruthy(test runtime.Value) bool {
	if runtime.Equals(test, nil) || runtime.Equals(test, runtime.Bool(false)) {
		return false
	}
	return true
}

func rootRT(rt runtime.Runtime) runtime.Runtime {
	for rt.Parent() != nil {
		rt = rt.Parent()
	}
	return rt
}
