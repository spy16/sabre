package core

import (
	"errors"
	"fmt"
	"strings"
)

// Nil represents a nil value.
type Nil struct{}

// Source returns the literal representation of nil value.
func (n Nil) Source() string { return "nil" }

// Bool represents a boolean value.
type Bool bool

// Source returns the literal representation of boolean value.
func (b Bool) Source() string { return fmt.Sprintf("%t", b) }

// Float64 represents double precision floating point numbers.
type Float64 float64

// Source returns the literal representation of floating point numbers.
func (f64 Float64) Source() string { return fmt.Sprintf("%f", f64) }

// Int64 represents integer values represented using decimal, octal, radix and
// hexadecimal formats.
type Int64 int64

// Source returns the literal representation of integer numbers.
func (i64 Int64) Source() string { return fmt.Sprintf("%d", i64) }

// Character represents a character literal.  For example, \a, \b, \1, \∂ etc
// are valid character literals. In addition, special literals like \newline,
// \space etc are supported by the reader.
type Character rune

// Source returns the literal representation of character (e.g., \a, \∂ etc.)
func (char Character) Source() string { return fmt.Sprintf("\\%c", rune(char)) }

// String represents double-quoted string literals. String Form represents
// the true string value obtained from the reader. Escape sequences are not
// applicable at this level.
type String string

// Source returns the double-quoted literal representation of string.
func (str String) Source() string { return fmt.Sprintf("\"%s\"", str) }

// Chars returns the string as a sequence of character values.
func (str String) Chars() Values {
	var vals Values
	for _, r := range str {
		vals = append(vals, Character(r))
	}
	return vals
}

// Keyword represents a keyword literal.
type Keyword string

// Source returns the literal representation of keyword (e.g, :specimen).
func (kw Keyword) Source() string { return fmt.Sprintf(":%s", kw) }

// Invoke enables keyword lookup for maps.
func (kw Keyword) Invoke(env Env, args ...Value) (Value, error) {
	argVals, err := EvalAll(env, args)
	if err != nil {
		return nil, err
	}

	if len(argVals) != 1 && len(argVals) != 2 {
		return nil, fmt.Errorf("invoke requires 1 or 2 args, got %d", len(argVals))
	}

	hm, ok := argVals[0].(*HashMap)
	if !ok {
		return Nil{}, nil
	}

	def := Value(Nil{})
	if len(argVals) == 2 {
		def = argVals[1]
	}

	return hm.Get(kw, def), nil
}

// Symbol represents a name given to a value in memory.
type Symbol struct {
	Position
	Value string
}

// Eval resolves the binding for the symbol and returns the associated value.
func (sym Symbol) Eval(env Env) (Value, error) {
	fields := strings.SplitN(sym.Value, ".", 2)

	if sym.Value == "." {
		fields = []string{"."}
	}

	target, err := env.Resolve(fields[0])
	if len(fields) == 1 || err != nil {
		return target, err
	}

	acc, ok := target.(memberAccessor)
	if !ok {
		return nil, errors.New("")
	}

	return acc.AccessMember(fields[1])
}

// Compare compares this symbol to the given value. Returns true if the given
// value is a symbol with same data.
func (sym Symbol) Compare(v Value) bool {
	other, ok := v.(Symbol)
	if !ok {
		return false
	}
	return other.Value == sym.Value
}

// Source returns the literal representation of the symbol which is the symbol
// value itself.
func (sym Symbol) Source() string { return sym.Value }

type memberAccessor interface {
	AccessMember(name string) (Value, error)
}
