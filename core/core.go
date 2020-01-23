package core

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/spy16/sabre"
)

// BindAll binds all core functions into the given scope.
func BindAll(scope sabre.Scope) error {
	core := map[string]sabre.Value{
		"λ":            sabre.GoFunc(Lambda),
		"fn":           sabre.GoFunc(Lambda),
		"do":           sabre.GoFunc(Do),
		"def":          sabre.GoFunc(Def),
		"eval":         sabre.GoFunc(Eval),
		"quote":        sabre.GoFunc(SimpleQuote),
		"syntax-quote": sabre.GoFunc(SyntaxQuote),
		"not":          Fn(Not),
		"error":        Fn(RaiseErr),
		"boolean":      Fn(MakeBool),
		"str":          Fn(MakeString),
		"type":         Fn(TypeOf),
		"set":          MakeContainer(sabre.Set(nil)),
		"list":         MakeContainer(sabre.List(nil)),
		"vector":       MakeContainer(sabre.Vector(nil)),
		"nil?":         IsType(reflect.TypeOf(sabre.Nil{})),
		"int?":         IsType(reflect.TypeOf(sabre.Int64(0))),
		"set?":         IsType(reflect.TypeOf(sabre.Set(nil))),
		"boolean?":     IsType(reflect.TypeOf(sabre.Bool(false))),
		"list?":        IsType(reflect.TypeOf(sabre.List(nil))),
		"string?":      IsType(reflect.TypeOf(sabre.String(""))),
		"float?":       IsType(reflect.TypeOf(sabre.Float64(0))),
		"vector?":      IsType(reflect.TypeOf(sabre.Vector(nil))),
		"keyword?":     IsType(reflect.TypeOf(sabre.Keyword(""))),
		"symbol?":      IsType(reflect.TypeOf(sabre.Symbol{})),
	}

	for sym, val := range core {
		if err := scope.Bind(sym, val); err != nil {
			return err
		}
	}

	return nil
}

// Eval evaluates the first argument and returns the result.
func Eval(scope sabre.Scope, args []sabre.Value) (sabre.Value, error) {
	vals, err := evalValueList(scope, args)
	if err != nil {
		return nil, err
	}

	if err := verifyArgCount([]int{1}, args); err != nil {
		return nil, err
	}

	return vals[0].Eval(scope)
}

// Lambda defines an anonymous function and returns.
func Lambda(scope sabre.Scope, args []sabre.Value) (sabre.Value, error) {
	if err := verifyArgCount([]int{1, 2}, args); err != nil {
		return nil, err
	}

	lArgs, isVector := args[0].(sabre.Vector)
	if !isVector {
		return nil, fmt.Errorf("first argument must be a vector of symbols")
	}

	lambdaBody := args[1:]
	lambdaArgs, err := toSymbolList(lArgs)
	if err != nil {
		return nil, err
	}

	return LambdaFn(scope, lambdaArgs, lambdaBody), nil
}

// LambdaFn creates a lambda function with given arguments and body.
func LambdaFn(scope sabre.Scope, argNames []sabre.Symbol, body []sabre.Value) sabre.GoFunc {
	return sabre.GoFunc(func(_ sabre.Scope, args []sabre.Value) (sabre.Value, error) {
		argVals, err := evalValueList(scope, args)
		if err != nil {
			return nil, err
		}

		if err := verifyArgCount([]int{len(argNames)}, argVals); err != nil {
			return nil, err
		}

		fnScope := sabre.NewScope(scope)
		for idx := range argNames {
			if err := fnScope.Bind(argNames[idx].String(), argVals[idx]); err != nil {
				return nil, err
			}
		}

		return sabre.Module(body).Eval(fnScope)
	})
}

// Def adds a binding to the scope. First argument must be a symbol
// and second argument must be a value.
func Def(scope sabre.Scope, args []sabre.Value) (sabre.Value, error) {
	if err := verifyArgCount([]int{2}, args); err != nil {
		return nil, err
	}

	sym, isSymbol := args[0].(sabre.Symbol)
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

	return sym, nil
}

// Not returns the negated version of the argument value.
func Not(args []sabre.Value) (sabre.Value, error) {
	if err := verifyArgCount([]int{1}, args); err != nil {
		return nil, err
	}

	return sabre.Bool(!isTruthy(args[0])), nil
}

// Do evaluates all the arguments and returns the result of last evaluation.
func Do(scope sabre.Scope, args []sabre.Value) (sabre.Value, error) {
	return sabre.Module(args).Eval(scope)
}

// RaiseErr signals an error. Stringified versions of args will be
// concatenated and used as error message.
func RaiseErr(vals []sabre.Value) (sabre.Value, error) {
	return nil, errors.New(string(stringFromVals(vals)))
}

// SimpleQuote prevents a form from being evaluated.
func SimpleQuote(scope sabre.Scope, forms []sabre.Value) (sabre.Value, error) {
	if err := verifyArgCount([]int{1}, forms); err != nil {
		return nil, err
	}

	return forms[0], nil
}

// SyntaxQuote recursively applies the quoting to the form.
func SyntaxQuote(scope sabre.Scope, forms []sabre.Value) (sabre.Value, error) {
	if err := verifyArgCount([]int{1}, forms); err != nil {
		return nil, err
	}

	quoteScope := sabre.NewScope(scope)
	quoteScope.Bind("unquote", sabre.GoFunc(unquote))

	return recursiveQuote(quoteScope, forms[0])
}

func unquote(scope sabre.Scope, forms []sabre.Value) (sabre.Value, error) {
	if err := verifyArgCount([]int{1}, forms); err != nil {
		return nil, err
	}

	return forms[0].Eval(scope)
}

func recursiveQuote(scope sabre.Scope, f sabre.Value) (sabre.Value, error) {
	switch v := f.(type) {
	case sabre.List:
		if isUnquote(v) {
			return f.Eval(scope)
		}

		quoted, err := quoteList(scope, v)
		return sabre.List(quoted), err

	case sabre.Set:
		quoted, err := quoteList(scope, v)
		return sabre.Set(quoted), err

	case sabre.Vector:
		quoted, err := quoteList(scope, v)
		return sabre.Vector(quoted), err

	default:
		return f, nil
	}
}

func isUnquote(list []sabre.Value) bool {
	if len(list) == 0 {
		return false
	}

	sym, isSymbol := list[0].(sabre.Symbol)
	if !isSymbol {
		return false
	}

	return sym.Value == "unquote"
}

func quoteList(scope sabre.Scope, forms []sabre.Value) ([]sabre.Value, error) {
	var quoted []sabre.Value
	for _, form := range forms {
		q, err := recursiveQuote(scope, form)
		if err != nil {
			return nil, err
		}

		quoted = append(quoted, q)
	}

	return quoted, nil
}
