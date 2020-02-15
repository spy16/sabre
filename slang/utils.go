package slang

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spy16/sabre"
)

func evalValueList(scope sabre.Scope, vals []sabre.Value) ([]sabre.Value, error) {
	var result []sabre.Value

	for _, arg := range vals {
		v, err := arg.Eval(scope)
		if err != nil {
			return nil, err
		}

		result = append(result, v)
	}

	return result, nil
}

func stringFromVals(vals []sabre.Value) sabre.String {
	argc := len(vals)
	switch argc {
	case 0:
		return sabre.String("")

	case 1:
		return sabre.String(strings.Trim(vals[0].String(), "\""))

	default:
		var sb strings.Builder
		for _, v := range vals {
			sb.WriteString(strings.Trim(v.String(), "\""))
		}
		return sabre.String(sb.String())
	}
}

func isTruthy(v sabre.Value) bool {
	var sabreNil = sabre.Nil{}
	if v == sabreNil {
		return false
	}

	if b, ok := v.(sabre.Bool); ok {
		return bool(b)
	}

	return true
}

func verifyArgCount(arities []int, args []sabre.Value) error {
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
