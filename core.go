package sabre

import (
	"fmt"
	"reflect"
)

// Def adds a binding to the scope. First argument must be a symbol
// and second argument must be a value.
func Def(scope Scope, args []Value) (Value, error) {
	if len(args) != 2 {
		return nil, arityErr(2, len(args), "")
	}

	sym, isSymbol := args[0].(Symbol)
	if !isSymbol {
		return nil, fmt.Errorf("first argument must be symbol, not '%v'", reflect.TypeOf(args[0]))
	}

	v, err := args[1].Eval(scope)
	if err != nil {
		return nil, err
	}

	if err := scope.Bind(sym.String(), v); err != nil {
		return nil, err
	}

	return List{Symbol("quote"), sym}, nil
}

// Lambda defines an anonymous function and returns.
func Lambda(_ Scope, args []Value) (Value, error) {
	if len(args) < 2 {
		return nil, arityErr(2, len(args), "")
	}

	lArgs, isVector := args[0].(Vector)
	if !isVector {
		return nil, fmt.Errorf("first argument must be a vector of symbols")
	}

	lambdaBody := args[1:]
	lambdaArgs, err := toSymbolList(lArgs)
	if err != nil {
		return nil, err
	}

	return LambdaFn(lambdaArgs, lambdaBody), nil
}

// LambdaFn creates a lambda function with given arguments and body.
func LambdaFn(argNames []Symbol, body []Value) Fn {
	return Fn(func(scope Scope, args []Value) (Value, error) {
		if len(args) != len(argNames) {
			return nil, arityErr(len(argNames), len(args), "")
		}

		fnScope := NewScope(scope, false)
		for idx := range argNames {
			if err := fnScope.Bind(argNames[idx].String(), args[idx]); err != nil {
				return nil, err
			}
		}

		return Module(body).Eval(fnScope)
	})
}

func evalFn(scope Scope, args []Value) (Value, error) {
	if len(args) != 1 {
		return nil, arityErr(1, len(args), "")
	}

	v, err := args[0].Eval(scope)
	if err != nil {
		return nil, err
	}

	return v.Eval(scope)
}

func toSymbolList(vals []Value) ([]Symbol, error) {
	var argNames []Symbol

	for _, arg := range vals {
		sym, isSymbol := arg.(Symbol)
		if !isSymbol {
			return nil, fmt.Errorf("first argument must be a vector of symbols")
		}

		argNames = append(argNames, sym)
	}

	return argNames, nil
}
