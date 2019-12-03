package sabre

import (
	"errors"
	"fmt"
	"strings"
)

// Form represents a LISP form.
type Form interface {
	Eval(scope Scope) Value
}

// Number represents a numerical value. If IsFloat is true, then value of
// the number will be in Float field. Otherwise,  value is in  Int. Three
// fields are  used to avoid  usage of interface{}, nil value  issues and
// also to keep allocations fixed.
type Number struct {
	IsFloat bool
	Float   float64
	Int     int64
}

// Eval returns the parsed numerical value of the literal.
func (ne Number) Eval(_ Scope) Value {
	if ne.IsFloat {
		return ne.Float
	}

	return ne.Int
}

func (ne Number) String() string {
	if ne.IsFloat {
		return fmt.Sprintf("%f", ne.Eval(nil))
	}

	return fmt.Sprintf("%d", ne.Eval(nil))
}

// String represents double-quoted string literals. String Form represents
// the true string value obtained from the reader. Escape sequences are not
// applicable at this level.
type String string

// Eval returns the unquoted string value of the literal.
func (se String) Eval(_ Scope) Value {
	return string(se)
}

func (se String) String() string {
	return fmt.Sprintf("\"%s\"", string(se))
}

// Symbol represents a name given to a value in memory.
type Symbol string

// Eval performs a lookup in the scope and returns the value bound to the
// ident. Returns error if lookup fails.
func (sym Symbol) Eval(scope Scope) Value {
	v, err := scope.Get(string(sym))
	if err != nil {
		return err
	}

	return v
}

// Character represents a character literal.  For example, \a, \b, \1, \∂ etc
// are valid character literals. In addition, special literals like \newline,
// \space etc are supported.
type Character rune

// Eval returns the character literal.
func (char Character) Eval(_ Scope) Value {
	return rune(char)
}

func (char Character) String() string {
	return fmt.Sprintf("\\%c", rune(char))
}

// Keyword represents a keyword literal.
type Keyword string

// Eval simply returns the keyword itself.
func (kw Keyword) Eval(_ Scope) Value {
	return string(kw)
}

// List represents an list of forms. Evaluating a list leads to a function
// invocation.
type List struct {
	Forms []Form
}

// Eval performs a function invocation by resolving first item in the list
// to a callable.
func (le List) Eval(scope Scope) Value {
	return errors.New("failed to invoke")
}

func (le List) String() string {
	parts := make([]string, len(le.Forms))
	for i, expr := range le.Forms {
		parts[i] = fmt.Sprintf("%s", expr)
	}
	return "(" + strings.Join(parts, " ") + ")"
}

// Vector represents a list of values. Unlike List type, evaluation of
// vector does not lead to function invoke.
type Vector []Form

// Eval simply returns the values that the vector holds.
func (vec Vector) Eval(scope Scope) Value {
	var vals []interface{}

	for _, expr := range vec {
		v := expr.Eval(scope)
		if e, ok := v.(error); ok {
			return e
		}

		vals = append(vals, v)
	}

	// TODO: Cache these vals?
	return vals
}

func (vec Vector) String() string {
	parts := make([]string, len(vec))
	for i, expr := range vec {
		parts[i] = fmt.Sprintf("%s", expr)
	}

	return "[" + strings.Join(parts, " ") + "]"
}

// Module represents a group of forms. Evaluating a module form returns the
// result of evaluating the last form in the list.
type Module []Form

// Eval evaluates every form in the module and returns the result of the last
// evaluation.
func (mod Module) Eval(scope Scope) Value {
	var v Value
	for _, expr := range mod {
		v = expr.Eval(scope)
		if _, ok := v.(error); ok {
			return v
		}
	}

	return v
}
