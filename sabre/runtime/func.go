package runtime

import (
	"fmt"
	"reflect"
)

// GoFunc provides a simple Go native function based invokable value.
type GoFunc func(rt Runtime, args ...Value) (Value, error)

// Equals returns true if the 'other' value is a GoFunc and has the same
// memory address (pointer value).
func (fn GoFunc) Equals(other Value) bool {
	gf, ok := other.(GoFunc)
	return ok && reflect.ValueOf(fn).Pointer() == reflect.ValueOf(gf).Pointer()
}

// Invoke simply dispatches the invocation request to the wrapped function.
// Wrapped function value receives un-evaluated list of arguments.
func (fn GoFunc) Invoke(env Runtime, args ...Value) (Value, error) {
	return fn(env, args...)
}

func (fn GoFunc) String() string { return fmt.Sprintf("GoFunc{%p}", fn) }
