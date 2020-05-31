package core

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

var _ Invokable = GoFunc(nil)

// Compare compares two values in an identity independent manner. If v1 implements
// `Comparable` interface then the comparison is delegated to the Compare() method.
func Compare(v1, v2 Value) bool {
	if isNil(v1) && isNil(v2) {
		return true
	}

	if cmp, ok := v1.(Comparable); ok {
		return cmp.Compare(v2)
	}

	return reflect.DeepEqual(v1, v2)
}

// EvalAll evaluates each value in the list against the given env and returns a list
// of resultant value.
func EvalAll(env Env, vals []Value) ([]Value, error) {
	var results []Value
	for _, f := range vals {
		res, err := env.Eval(f)
		if err != nil {
			return nil, err
		}
		results = append(results, res)
	}
	return results, nil
}

// SeqString returns a string representation for the sequence with given prefix
// suffix and separator.
func SeqString(seq Seq, begin, end, sep string) string {
	var parts []string
	ForEach(seq, func(item Value) bool {
		parts = append(parts, item.String())
		return false
	})
	return begin + strings.Join(parts, sep) + end
}

// ForEach reads from the sequence and calls the given function for each item.
// Function can return true to stop the iteration.
func ForEach(seq Seq, call func(item Value) bool) {
	for seq != nil {
		v := seq.First()
		if v == nil || call(seq.First()) {
			break
		}
		seq = seq.Next()
	}
}

// VerifyArgCount checks the arg count against the given possible arities and
// returns clean errors with appropriate hints if the arg count doesn't match
// any arity.
func VerifyArgCount(arities []int, argCount int) error {
	sort.Ints(arities)

	if len(arities) == 0 && argCount != 0 {
		return fmt.Errorf("call requires no arguments, got %d", argCount)
	}

	switch len(arities) {
	case 1:
		if argCount != arities[0] {
			return fmt.Errorf("call requires exactly %d argument(s), got %d", arities[0], argCount)
		}

	case 2:
		c1, c2 := arities[0], arities[1]
		if argCount != c1 && argCount != c2 {
			return fmt.Errorf("call requires %d or %d argument(s), got %d", c1, c2, argCount)
		}

	default:
		return fmt.Errorf("wrong number of arguments (%d) passed", argCount)
	}

	return nil
}

func isNil(v Value) bool {
	_, isNil := v.(Nil)
	return v == nil || isNil
}

// GoFunc provides a simple Go native function based invokable value.
type GoFunc func(env Env, args ...Value) (Value, error)

// Eval simply returns itself.
func (fn GoFunc) Eval(_ Env) (Value, error) { return fn, nil }

func (fn GoFunc) String() string {
	return fmt.Sprintf("GoFunc{}")
}

// Invoke simply dispatches the invocation request to the wrapped function.
func (fn GoFunc) Invoke(env Env, args ...Value) (Value, error) {
	return fn(env, args...)
}
