package sabre

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

func bindCore(scope Scope) error {
	core := map[string]Value{
		"λ":        Fn(Lambda),
		"fn":       Fn(Lambda),
		"do":       Fn(Do),
		"not":      Fn(Not),
		"def":      Fn(Def),
		"eval":     Fn(evalFn),
		"error":    Fn(RaiseErr),
		"boolean":  Fn(MakeBool),
		"str":      Fn(MakeString),
		"type":     Fn(TypeOf),
		"set":      MakeContainer(Set(nil)),
		"list":     MakeContainer(List(nil)),
		"vector":   MakeContainer(Vector(nil)),
		"nil?":     IsType(reflect.TypeOf(Nil{})),
		"int?":     IsType(reflect.TypeOf(Int64(0))),
		"set?":     IsType(reflect.TypeOf(Set(nil))),
		"boolean?": IsType(reflect.TypeOf(Bool(false))),
		"list?":    IsType(reflect.TypeOf(List(nil))),
		"string?":  IsType(reflect.TypeOf(String(""))),
		"float?":   IsType(reflect.TypeOf(Float64(0))),
		"vector?":  IsType(reflect.TypeOf(Vector(nil))),
		"symbol?":  IsType(reflect.TypeOf(Symbol(""))),
		"keyword?": IsType(reflect.TypeOf(Keyword(""))),
	}

	for sym, val := range core {
		if err := scope.Bind(sym, val); err != nil {
			return err
		}
	}

	return nil
}

// TypeOf returns the type information object for the given argument.
func TypeOf(scope Scope, args []Value) (Value, error) {
	argVals, err := evalValueList(scope, args)
	if err != nil {
		return nil, err
	}

	if err := verifyArgCount([]int{1}, argVals); err != nil {
		return nil, err
	}

	return Type{rt: reflect.TypeOf(argVals[0])}, nil
}

// RaiseErr signals an error. Stringified versions of args will be
// concatenated and used as error message.
func RaiseErr(scope Scope, args []Value) (Value, error) {
	argVals, err := evalValueList(scope, args)
	if err != nil {
		return nil, err
	}

	return nil, errors.New(string(stringFromVals(argVals)))
}

// MakeBool converts given argument to a boolean. Any truthy value
// is converted to true and else false.
func MakeBool(scope Scope, args []Value) (Value, error) {
	argVals, err := evalValueList(scope, args)
	if err != nil {
		return nil, err
	}

	if err := verifyArgCount([]int{1}, argVals); err != nil {
		return nil, err
	}

	return Bool(isTruthy(argVals[0])), nil
}

// MakeString returns stringified version of all args.
func MakeString(scope Scope, args []Value) (Value, error) {
	argVals, err := evalValueList(scope, args)
	if err != nil {
		return nil, err
	}

	return stringFromVals(argVals), nil
}

func stringFromVals(vals []Value) String {
	argc := len(vals)
	switch argc {
	case 0:
		return String("")

	case 1:
		return String(strings.Trim(vals[0].String(), "\""))

	default:
		var sb strings.Builder
		for _, v := range vals {
			sb.WriteString(strings.Trim(v.String(), "\""))
		}
		return String(sb.String())
	}
}

// MakeContainer can make a composite type like list, set and vector from
// given args.
func MakeContainer(targetType Value) Fn {
	return func(scope Scope, args []Value) (Value, error) {
		argVals, err := evalValueList(scope, args)
		if err != nil {
			return nil, err
		}

		switch targetType.(type) {
		case List:
			return List(argVals), nil

		case Vector:
			return Vector(argVals), nil

		case Set:
			return Set(argVals), nil
		}

		return nil, fmt.Errorf("cannot make container of type '%s'", reflect.TypeOf(targetType))
	}
}

// IsType returns a Fn that checks if the value is of given type.
func IsType(rt reflect.Type) Fn {
	return func(scope Scope, args []Value) (Value, error) {
		if err := verifyArgCount([]int{1}, args); err != nil {
			return nil, err
		}

		v, err := args[0].Eval(scope)
		if err != nil {
			return nil, err
		}

		target := reflect.TypeOf(v)
		return Bool(target == rt), nil
	}
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
func Lambda(scope Scope, args []Value) (Value, error) {
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

	return LambdaFn(scope, lambdaArgs, lambdaBody), nil
}

// LambdaFn creates a lambda function with given arguments and body.
func LambdaFn(scope Scope, argNames []Symbol, body []Value) Fn {
	return Fn(func(_ Scope, args []Value) (Value, error) {
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
