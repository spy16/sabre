package runtime

import (
	"fmt"
	"reflect"
	"strings"
)

var (
	_ Value     = (*Fn)(nil)
	_ Invokable = (*Fn)(nil)
)

// Fn represents a function definition.
type Fn struct {
	Name     string
	Args     []string
	Variadic bool
	Body     Value
}

// Invoke executes the function with given arguments.
func (fn *Fn) Invoke(rt Runtime, args ...Value) (Value, error) {
	fnEnv := New(rt)

	for idx := range fn.Args {
		var argVal Value
		if idx == len(fn.Args)-1 && fn.Variadic {
			argVal = NewSeq(args[idx:]...)
		} else {
			argVal = args[idx]
		}

		_ = fnEnv.Bind(fn.Args[idx], argVal)
	}

	if fn.Body == nil {
		return Nil{}, nil
	}

	return fnEnv.Eval(fn.Body)
}

// Eval returns the function itself.
func (fn *Fn) Eval(_ Runtime) (Value, error) { return fn, nil }

func (fn Fn) String() string { return fn.stringWithPrefix("fn") }

// Equals returns true if 'other' is also an Fn value and has the same
// signature and body.
func (fn *Fn) Equals(other Value) bool {
	otherFn, ok := other.(*Fn)
	if !ok || otherFn == nil {
		return false
	}

	sameArgs := reflect.DeepEqual(fn.Args, otherFn.Args)
	bothVariadic := fn.Variadic == otherFn.Variadic
	return bothVariadic && sameArgs && Equals(fn.Body, otherFn.Body)
}

func (fn Fn) stringWithPrefix(prefix string) string {
	var sb strings.Builder
	sb.WriteRune('(')

	if prefix != "" {
		sb.WriteString(prefix + " ")
	}

	sb.WriteRune('[')
	if len(fn.Args) > 0 {
		if fn.Variadic {
			sb.WriteString(strings.Join(fn.Args[:len(fn.Args)-1], " "))
			sb.WriteString(" & " + fn.Args[len(fn.Args)-1])
		} else {
			sb.WriteString(strings.Join(fn.Args, " "))
		}
	}
	sb.WriteString("] ")
	if fn.Body != nil {
		sb.WriteString(fn.Body.String())
	}
	sb.WriteRune(')')

	return sb.String()
}

// GoFunc provides a simple Go native function based invokable value.
type GoFunc func(rt Runtime, args ...Value) (Value, error)

// Eval simply returns itself.
func (fn GoFunc) Eval(_ Runtime) (Value, error) { return fn, nil }

// Equals returns true if the 'other' value is a GoFunc and has the same
// memory address (pointer value).
func (fn GoFunc) Equals(other Value) bool {
	gf, ok := other.(GoFunc)
	return ok && reflect.ValueOf(fn).Pointer() == reflect.ValueOf(gf).Pointer()
}

func (fn GoFunc) String() string {
	return fmt.Sprintf("GoFunc{%p}", fn)
}

// Invoke simply dispatches the invocation request to the wrapped function.
// Wrapped function value receives un-evaluated list of arguments.
func (fn GoFunc) Invoke(env Runtime, args ...Value) (Value, error) {
	return fn(env, args...)
}
