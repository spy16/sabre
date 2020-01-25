package sabre

import (
	"fmt"
	"reflect"
	"sort"
)

// New returns an instance of MapScope with all the special forms setup.
func New() *MapScope {
	specials := map[string]Value{
		"λ":            GoFunc(Lambda),
		"fn":           GoFunc(Lambda),
		"do":           GoFunc(Do),
		"def":          GoFunc(Def),
		"if":           GoFunc(If),
		"quote":        GoFunc(SimpleQuote),
		"syntax-quote": GoFunc(SyntaxQuote),
	}

	scope := NewScope(nil)
	for name, val := range specials {
		_ = scope.Bind(name, val)
	}
	return scope
}

// If implments if-conditional flow using (if test then else?) form.
func If(scope Scope, args []Value) (Value, error) {
	if err := verifyArgCount([]int{2, 3}, args); err != nil {
		return nil, err
	}

	test, err := args[0].Eval(scope)
	if err != nil {
		return nil, err
	}

	if !isTruthy(test) {
		// handle 'else' flow.
		if len(args) == 2 {
			return Nil{}, nil
		}

		return args[2].Eval(scope)
	}

	// handle 'if true' flow.
	return args[1].Eval(scope)
}

// Def adds a binding to the scope. Should have the form (def symbol value).
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

	return sym, nil
}

// Do evaluates all the arguments and returns the result of last evaluation.
// Must have the form (do <expr>*)
func Do(scope Scope, args []Value) (Value, error) {
	return Module(args).Eval(scope)
}

// Lambda defines an anonymous function and returns. Must have the form
// (fn [arg*] expr*)
func Lambda(scope Scope, args []Value) (Value, error) {
	if err := verifyArgCount([]int{1, 2}, args); err != nil {
		return nil, err
	}

	lArgs, isVector := args[0].(Vector)
	if !isVector {
		return nil, fmt.Errorf("first argument must be a vector of symbols")
	}

	lambdaBody := args[1:]

	lambdaArgs, err := toSymbolList(lArgs.Items)
	if err != nil {
		return nil, err
	}

	return LambdaFn(scope, lambdaArgs, lambdaBody), nil
}

// LambdaFn creates a lambda function with given arguments and body.
func LambdaFn(scope Scope, argNames []Symbol, body []Value) GoFunc {
	return GoFunc(func(_ Scope, args []Value) (Value, error) {
		argVals, err := evalValueList(scope, args)
		if err != nil {
			return nil, err
		}

		if err := verifyArgCount([]int{len(argNames)}, argVals); err != nil {
			return nil, err
		}

		fnScope := NewScope(scope)
		for idx := range argNames {
			if err := fnScope.Bind(argNames[idx].String(), argVals[idx]); err != nil {
				return nil, err
			}
		}

		return Module(body).Eval(fnScope)
	})
}

// SimpleQuote prevents a form from being evaluated.
func SimpleQuote(scope Scope, forms []Value) (Value, error) {
	if err := verifyArgCount([]int{1}, forms); err != nil {
		return nil, err
	}

	return forms[0], nil
}

// SyntaxQuote recursively applies the quoting to the form.
func SyntaxQuote(scope Scope, forms []Value) (Value, error) {
	if err := verifyArgCount([]int{1}, forms); err != nil {
		return nil, err
	}

	quoteScope := NewScope(scope)
	quoteScope.Bind("unquote", GoFunc(unquote))

	return recursiveQuote(quoteScope, forms[0])
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

func unquote(scope Scope, forms []Value) (Value, error) {
	if err := verifyArgCount([]int{1}, forms); err != nil {
		return nil, err
	}

	return forms[0].Eval(scope)
}

func recursiveQuote(scope Scope, f Value) (Value, error) {
	switch v := f.(type) {
	case List:
		if isUnquote(v.Items) {
			return f.Eval(scope)
		}

		quoted, err := quoteList(scope, v.Items)
		return List{Items: quoted}, err

	case Set:
		quoted, err := quoteList(scope, v.Items)
		return Set{Items: quoted}, err

	case Vector:
		quoted, err := quoteList(scope, v.Items)
		return Vector{Items: quoted}, err

	default:
		return f, nil
	}
}

func isUnquote(list []Value) bool {
	if len(list) == 0 {
		return false
	}

	sym, isSymbol := list[0].(Symbol)
	if !isSymbol {
		return false
	}

	return sym.Value == "unquote"
}

func quoteList(scope Scope, forms []Value) ([]Value, error) {
	var quoted []Value
	for _, form := range forms {
		q, err := recursiveQuote(scope, form)
		if err != nil {
			return nil, err
		}

		quoted = append(quoted, q)
	}

	return quoted, nil
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
