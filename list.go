package sabre

import (
	"fmt"
	"reflect"
)

// List represents an list of forms/vals. Evaluating a list leads to a
// function invocation.
type List []Value

// Eval performs an invocation.
func (lf List) Eval(scope Scope) (Value, error) {
	if len(lf) == 0 {
		return List(nil), nil
	}

	if isQuote(lf[0]) {
		return quoteValue(lf[1:])
	}

	target, err := lf[0].Eval(scope)
	if err != nil {
		return nil, err
	}

	if specialFn, ok := target.(SpecialFn); ok {
		return specialFn(scope, lf[1:])
	}

	fn, ok := target.(invokable)
	if !ok {
		return nil, fmt.Errorf("cannot invoke value of type '%s'", reflect.TypeOf(lf[0]))
	}

	argVals, err := evalValueList(scope, lf[1:])
	if err != nil {
		return nil, err
	}

	return fn.Invoke(argVals)
}

func (lf List) String() string {
	return containerString(lf, "(", ")", " ")
}

func isQuote(v Value) bool {
	sym, isSymbol := v.(Symbol)
	if !isSymbol {
		return false
	}

	return sym == "quote"
}

func quoteValue(args []Value) (Value, error) {
	if len(args) != 1 {
		return nil, arityErr(1, len(args), "")
	}

	return args[0], nil
}

func evalValueList(scope Scope, vals []Value) ([]Value, error) {
	var result []Value

	for _, arg := range vals {
		v, err := arg.Eval(scope)
		if err != nil {
			return nil, err
		}

		result = append(result, v)
	}

	return result, nil
}

func readList(rd *Reader, _ rune) (Value, error) {
	forms, err := readContainer(rd, '(', ')', "list")
	if err != nil {
		return nil, err
	}

	return List(forms), nil
}

type invokable interface {
	Invoke(argVals []Value) (Value, error)
}
