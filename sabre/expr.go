package sabre

import (
	"fmt"
	"reflect"
	"strings"
)

var (
	_ Expr = (*DefExpr)(nil)
	_ Expr = (*ConstExpr)(nil)
	_ Expr = (*InvokeExpr)(nil)
)

// Expr represents an expression that can be evaluated against Sabre instance.
// Evaluation might have side effects on the instance.
type Expr interface {
	// Eval evaluates the form against the environment and returns the result
	// or error.
	Eval(env *Sabre) (Value, error)
}

// ConstExpr simply returns the value wrapped inside and has no side-effect
// on the environment.
type ConstExpr struct{ Value Value }

// Eval simply returns the wrapped value. Eval has no side-effect on the env.
func (ce *ConstExpr) Eval(_ *Sabre) (Value, error) { return ce.Value, nil }

// DefExpr creates global bindings in the given environment when evaluated.
type DefExpr struct {
	Name  string
	Value Value
}

// Eval validates and adds the binding to the global frame of env.
func (de *DefExpr) Eval(env *Sabre) (Value, error) {
	de.Name = strings.TrimSpace(de.Name)
	if de.Name == "" {
		return nil, fmt.Errorf("invalid name for bind: '%s'", de.Name)
	}
	env.stack[0].vars[de.Name] = de.Value
	return Symbol{Value: de.Name}, nil
}

// InvokeExpr performs invocation of target when evaluated.
type InvokeExpr struct {
	Target Expr
	Args   []Expr
}

// Eval evaluates target and argument exprs and invokes the target result with
// arg results. Returns error if target result is not Invokable.
func (ie *InvokeExpr) Eval(s *Sabre) (Value, error) {
	val, err := ie.Target.Eval(s)
	if err != nil {
		return nil, err
	}

	fn, ok := val.(Invokable)
	if !ok {
		return nil, fmt.Errorf("value of type '%s' is not invokable",
			reflect.TypeOf(val))
	}

	var args []Value
	for _, ae := range ie.Args {
		v, err := ae.Eval(s)
		if err != nil {
			return nil, err
		}
		args = append(args, v)
	}

	return fn.Invoke(s, args...)
}
