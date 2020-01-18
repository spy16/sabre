package sabre

import (
	"fmt"
	"reflect"
	"sort"
)

func bindCore(scope Scope) error {
	core := map[string]Value{
		"λ":    Fn(Lambda),
		"fn":   Fn(Lambda),
		"do":   Fn(Do),
		"not":  Fn(Not),
		"def":  Fn(Def),
		"eval": Fn(evalFn),
	}

	for sym, val := range core {
		if err := scope.Bind(sym, val); err != nil {
			return err
		}
	}

	return nil
}

// Do evaluates all the arguments and returns the result of last evaluation.
func Do(scope Scope, args []Value) (Value, error) {
	return Module(args).Eval(scope)
}

// Not returns the negated version of the argument value.
func Not(scope Scope, args []Value) (Value, error) {
	if err := verifyArgCount([]int{1}, args); err != nil {
		return nil, err
	}

	v, err := args[0].Eval(scope)
	if err != nil {
		return nil, err
	}

	return Bool(!isTruthy(v)), nil
}

// Def adds a binding to the scope. First argument must be a symbol
// and second argument must be a value.
func Def(scope Scope, args []Value) (Value, error) {
	if err := verifyArgCount([]int{2}, args); err != nil {
		return nil, err
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
	if err := verifyArgCount([]int{1, 2}, args); err != nil {
		return nil, err
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
		argVals, err := evalValueList(scope, args)
		if err != nil {
			return nil, err
		}

		if err := verifyArgCount([]int{len(argNames)}, argVals); err != nil {
			return nil, err
		}

		fnScope := NewScope(scope, false)
		for idx := range argNames {
			if err := fnScope.Bind(argNames[idx].String(), argVals[idx]); err != nil {
				return nil, err
			}
		}

		return Module(body).Eval(fnScope)
	})
}

func isTruthy(v Value) bool {
	var sabreNil = Nil{}
	if v == sabreNil {
		return false
	}

	if b, ok := v.(Bool); ok {
		return bool(b)
	}

	return true
}

func evalFn(scope Scope, args []Value) (Value, error) {
	if err := verifyArgCount([]int{1}, args); err != nil {
		return nil, err
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

func verifyArgCount(arities []int, args []Value) error {
	actual := len(args)
	sort.Ints(arities)

	if len(arities) == 0 && actual != 0 {
		return fmt.Errorf("call requires no arguments, got %d", actual)
	}

	L := len(arities)
	switch {
	case L == 1 && actual != arities[0]:
		return fmt.Errorf("call requires exactly %d argument(s), got %d", arities[0], actual)

	case L == 2:
		c1, c2 := arities[0], arities[1]
		if actual != c1 && actual != c2 {
			return fmt.Errorf("call requires %d or %d argument(s), got %d", c1, c2, actual)
		}

	case L > 2:
		return fmt.Errorf("wrong number of arguments (%d) passed", actual)
	}

	return nil
}
