package slang

import (
	"reflect"

	"github.com/spy16/sabre"
)

// BindAll binds all core functions into the given scope.
func BindAll(scope sabre.Scope) error {
	core := map[string]sabre.Value{
		"eval":     sabre.GoFunc(Eval),
		"not":      Fn(Not),
		"boolean":  Fn(MakeBool),
		"str":      Fn(MakeString),
		"type":     Fn(TypeOf),
		"set":      makeContainer(sabre.Set{}),
		"list":     makeContainer(&sabre.List{}),
		"vector":   makeContainer(sabre.Vector{}),
		"nil?":     IsType(reflect.TypeOf(sabre.Nil{})),
		"int?":     IsType(reflect.TypeOf(sabre.Int64(0))),
		"set?":     IsType(reflect.TypeOf(sabre.Set{})),
		"boolean?": IsType(reflect.TypeOf(sabre.Bool(false))),
		"list?":    IsType(reflect.TypeOf(sabre.List{})),
		"string?":  IsType(reflect.TypeOf(sabre.String(""))),
		"float?":   IsType(reflect.TypeOf(sabre.Float64(0))),
		"vector?":  IsType(reflect.TypeOf(sabre.Vector{})),
		"keyword?": IsType(reflect.TypeOf(sabre.Keyword(""))),
		"symbol?":  IsType(reflect.TypeOf(sabre.Symbol{})),
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

// Not returns the negated version of the argument value.
func Not(args []sabre.Value) (sabre.Value, error) {
	if err := verifyArgCount([]int{1}, args); err != nil {
		return nil, err
	}

	return sabre.Bool(!isTruthy(args[0])), nil
}
