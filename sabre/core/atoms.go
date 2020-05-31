package core

import (
	"fmt"
)

var (
	_ Value = Nil{}
	_ Value = Bool(true)
	_ Value = Int64(0)
	_ Value = Float64(0)
	_ Value = Char('a')
	_ Value = Keyword("specimen")
	_ Value = String("specimen")
	_ Value = Symbol{}

	_ Invokable  = Keyword("specimen")
	_ Comparable = Symbol{}
)

// Nil represents a nil value.
type Nil struct{}

// Eval returns the underlying value.
func (n Nil) Eval(_ Env) (Value, error) { return n, nil }

func (n Nil) String() string { return "nil" }

// Bool represents a boolean value.
type Bool bool

// Eval returns the underlying value.
func (b Bool) Eval(_ Env) (Value, error) { return b, nil }

func (b Bool) String() string { return fmt.Sprintf("%t", b) }

// Float64 represents double precision floating point numbers represented
// using decimal or scientific number formats.
type Float64 float64

// Eval simply returns itself since Floats evaluate to themselves.
func (f64 Float64) Eval(_ Env) (Value, error) { return f64, nil }

func (f64 Float64) String() string { return fmt.Sprintf("%f", f64) }

// Int64 represents integer values represented using decimal, octal, radix
// and hexadecimal formats.
type Int64 int64

// Eval simply returns itself since Integers evaluate to themselves.
func (i64 Int64) Eval(_ Env) (Value, error) { return i64, nil }

func (i64 Int64) String() string { return fmt.Sprintf("%d", i64) }

// Char represents a character literal.  For example, \a, \b, \1, \∂ etc are
// valid character literals. In addition, special literals like \newline, \space
// etc are supported by the reader.
type Char rune

// Eval simply returns itself since Chracters evaluate to themselves.
func (char Char) Eval(_ Env) (Value, error) { return char, nil }

func (char Char) String() string { return fmt.Sprintf("\\%c", rune(char)) }

// String represents double-quoted string literals. String Form represents
// the true string value obtained from the reader. Escape sequences are not
// applicable at this level.
type String string

// Eval simply returns itself since Strings evaluate to themselves.
func (se String) Eval(_ Env) (Value, error) { return se, nil }

func (se String) String() string { return fmt.Sprintf("\"%s\"", string(se)) }

// Keyword represents a keyword literal.
type Keyword string

// Eval simply returns itself since Keywords evaluate to themselves.
func (kw Keyword) Eval(_ Env) (Value, error) { return kw, nil }

func (kw Keyword) String() string { return fmt.Sprintf(":%s", string(kw)) }

// Invoke enables keyword lookup for maps.
func (kw Keyword) Invoke(scope Env, args ...Value) (Value, error) {
	if err := VerifyArgCount([]int{1, 2}, len(args)); err != nil {
		return nil, err
	}

	argVals, err := EvalAll(scope, args)
	if err != nil {
		return nil, err
	}

	assocVal, ok := argVals[0].(Map)
	if !ok {
		return Nil{}, nil
	}

	def := Value(Nil{})
	if len(argVals) == 2 {
		def = argVals[1]
	}

	val, err := assocVal.Get(kw)
	if err != nil || val == nil {
		return def, err
	}

	return val, nil
}

// Symbol represents a name given to a value in memory.
type Symbol struct {
	Position
	Value string
}

// Eval returns the value bound to this symbol in current context.
func (sym Symbol) Eval(scope Env) (Value, error) {
	return scope.Resolve(sym.Value)
}

// Compare compares this symbol to the given value. Returns true if
// the given value is a symbol with same data.
func (sym Symbol) Compare(v Value) bool {
	other, ok := v.(Symbol)
	if !ok {
		return false
	}

	return other.Value == sym.Value
}

func (sym Symbol) String() string { return sym.Value }
