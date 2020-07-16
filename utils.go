package sabre

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spy16/sabre/runtime"
)

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
			return fmt.Errorf(
				"call requires exactly %d argument(s), got %d",
				arities[0], argCount,
			)
		}

	case 2:
		c1, c2 := arities[0], arities[1]
		if argCount != c1 && argCount != c2 {
			return fmt.Errorf(
				"call requires %d or %d argument(s), got %d", c1, c2, argCount)
		}

	default:
		for i := 0; i < len(arities); i++ {
			if arities[i] == argCount {
				return nil
			}
		}
		return fmt.Errorf("wrong number of arguments (%d) passed", argCount)
	}

	return nil
}

// Module is a convenience wrapper for evaluating multiple values. Evaluation
// returns last result of evaluating all forms in the module.
type Module []runtime.Value

// Eval evaluates each form in the module and returns the result of last eval.
// Returns Nil if the module is empty.
func (mod Module) Eval(rt runtime.Runtime) (runtime.Value, error) {
	var res runtime.Value
	var err error
	for _, f := range mod {
		res, err = rt.Eval(f)
		if err != nil {
			return nil, err
		}
	}

	if runtime.Equals(res, nil) {
		return runtime.Nil{}, nil
	}

	return res, nil
}

// String returns module represented using a (do expr*) form.
func (mod Module) String() string {
	var sb strings.Builder
	sb.WriteString("(do ")
	for i, f := range mod {
		sb.WriteString(f.String())
		if i < len(mod)-1 {
			sb.WriteString("\n    ")
		}
	}
	sb.WriteRune(')')
	return sb.String()
}
