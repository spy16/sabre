package runtime

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
)

var (
	_ Value     = GoFunc(nil)
	_ Invokable = GoFunc(nil)
)

// New returns an empty runtime with given parent runtime. Returned runtime does not
// support qualified symbol resolution. parent argument can be nil to make this the
// root runtime.
func New(parent Runtime) Runtime {
	return &mapEnv{
		scope:  map[string]Value{},
		parent: parent,
	}
}

// Equals compares two values in an identity independent manner. If v1 implements
// `Equals(Value)` method, then the comparison is delegated to it.
func Equals(v1, v2 Value) bool {
	if isNil(v1) && isNil(v2) {
		return true
	}

	if cmp, ok := v1.(interface{ Equals(other Value) bool }); ok {
		return cmp.Equals(v2)
	}

	s1, isV1Seq := v1.(Seq)
	s2, isV2Seq := v2.(Seq)
	if isV1Seq && isV2Seq {
		return compareSeq(s1, s2)
	}

	return reflect.DeepEqual(v1, v2)
}

// EvalAll evaluates each value in the list against the given env and returns a list
// of resultant value.
func EvalAll(rt Runtime, vals []Value) ([]Value, error) {
	var results []Value
	for _, f := range vals {
		res, err := rt.Eval(f)
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

// GoFunc provides a simple Go native function based invokable value.
type GoFunc func(env Runtime, args ...Value) (Value, error)

// Eval simply returns itself.
func (fn GoFunc) Eval(_ Runtime) (Value, error) { return fn, nil }

// Equals always returns false.
func (fn GoFunc) Equals(other Value) bool { return false }

func (fn GoFunc) String() string { return fmt.Sprintf("GoFunc{}") }

// Invoke simply dispatches the invocation request to the wrapped function.
func (fn GoFunc) Invoke(env Runtime, args ...Value) (Value, error) {
	return fn(env, args...)
}

func compareSeq(s1, s2 Seq) bool {
	if s1.Count() != s2.Count() {
		return false
	}

	for s1 != nil && s2 != nil {
		if !Equals(s1.First(), s2.First()) {
			return false
		}
		s1 = s1.Next()
		s2 = s2.Next()
	}

	return true
}

func isNil(v Value) bool {
	_, isNil := v.(Nil)
	return v == nil || isNil
}

type mapEnv struct {
	mu     sync.RWMutex
	scope  map[string]Value
	parent Runtime
}

func (env *mapEnv) Eval(form Value) (Value, error) {
	if isNil(form) {
		return Nil{}, nil
	}

	v, err := form.Eval(env)
	if err != nil {
		e := NewErr(false, getPosition(form), err)
		e.Form = form
		return nil, e
	}

	if v == nil {
		return Nil{}, nil
	}

	return v, nil
}

func (env *mapEnv) Bind(symbol string, v Value) error {
	env.mu.Lock()
	defer env.mu.Unlock()

	env.scope[symbol] = v
	return nil
}

func (env *mapEnv) Resolve(symbol string) (Value, error) {
	env.mu.RLock()
	defer env.mu.RUnlock()

	v, found := env.scope[symbol]
	if !found {
		if env.parent == nil {
			return nil, ErrNotFound
		}

		return env.parent.Resolve(symbol)
	}

	return v, nil
}

func (env *mapEnv) Parent() Runtime { return env.parent }

func getPosition(form Value) Position {
	return Position{}
}
