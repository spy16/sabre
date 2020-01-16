package sabre

import (
	"fmt"
)

// Bool represents a boolean value.
type Bool bool

// Eval returns the underlying value.
func (b Bool) Eval(_ Scope) (Value, error) { return b, nil }

func (b Bool) String() string { return fmt.Sprintf("%t", b) }

// Float64 represents double precision floating point numbers represented
// using float or scientific number formats.
type Float64 float64

// Eval returns the underlying value.
func (f64 Float64) Eval(_ Scope) (Value, error) { return f64, nil }

func (f64 Float64) String() string { return fmt.Sprintf("%f", f64) }

// Int64 represents integer values represented using decimal, octal, radix
// and hexadecimal formats.
type Int64 int64

// Eval returns the underlying value.
func (i64 Int64) Eval(_ Scope) (Value, error) { return i64, nil }

func (i64 Int64) String() string { return fmt.Sprintf("%d", i64) }

// String represents double-quoted string literals. String Form represents
// the true string value obtained from the reader. Escape sequences are not
// applicable at this level.
type String string

// Eval returns the underlying value.
func (se String) Eval(_ Scope) (Value, error) { return se, nil }

func (se String) String() string { return fmt.Sprintf("\"%s\"", string(se)) }

// Character represents a character literal.  For example, \a, \b, \1, \∂ etc
// are valid character literals. In addition, special literals like \newline,
// \space etc are supported.
type Character rune

// Eval returns the underlying value.
func (char Character) Eval(_ Scope) (Value, error) { return char, nil }

func (char Character) String() string { return fmt.Sprintf("\\%c", rune(char)) }

// Keyword represents a keyword literal.
type Keyword string

// Eval returns the underlying value.
func (kw Keyword) Eval(_ Scope) (Value, error) { return kw, nil }

func (kw Keyword) String() string { return fmt.Sprintf(":%s", string(kw)) }

// Symbol represents a name given to a value in memory.
type Symbol string

// Eval returns the underlying value.
func (sym Symbol) Eval(scope Scope) (Value, error) { return scope.Resolve(string(sym)) }

func (sym Symbol) String() string { return string(sym) }
