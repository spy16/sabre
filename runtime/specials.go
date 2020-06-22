package runtime

import (
	"errors"
	"fmt"
	"reflect"
)

func condForm(rt Runtime, args ...Value) (specialInvoke, error) {
	if len(args) > 0 && len(args)%2 != 0 {
		return nil, errors.New("cond requires even number of forms")
	}

	return func() (Value, error) {
		for i := 0; i < len(args); i += 2 {
			test, err := rt.Eval(args[i])
			if err != nil {
				return nil, err
			}

			if isTruthy(test) {
				return rt.Eval(args[i+1])
			}
		}
		return Nil{}, nil
	}, nil
}

func doForm(rt Runtime, args ...Value) (specialInvoke, error) {
	return func() (Value, error) {
		if len(args) == 0 {
			return Nil{}, nil
		}
		results, err := EvalAll(rt, args)
		if err != nil {
			return nil, err
		}
		return results[len(results)-1], nil
	}, nil
}

func quoteForm(_ Runtime, args ...Value) (specialInvoke, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("quote requires exactly 1 arg, got %d", len(args))
	}
	return func() (Value, error) { return args[0], nil }, nil
}

func defForm(rt Runtime, args ...Value) (specialInvoke, error) {
	if len(args) != 2 {
		return nil, errors.New("def requires exactly 2 arguments")
	}

	sym, isSym := args[0].(Symbol)
	if !isSym {
		return nil, fmt.Errorf("first arg to 'def' must be symbol, not %v",
			reflect.TypeOf(args[0]))
	}

	return func() (Value, error) {
		val, err := rt.Eval(args[1])
		if err != nil {
			return nil, err
		}
		return sym, rootRuntime(rt).Bind(sym.Value, val)
	}, nil
}

func rootRuntime(rt Runtime) Runtime {
	for rt.Parent() != nil {
		rt = rt.Parent()
	}
	return rt
}

func isTruthy(v Value) bool {
	if isNil(v) {
		return false
	}
	boolVal, isBool := v.(Bool)
	return !isBool || bool(boolVal)
}
