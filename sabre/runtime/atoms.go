package runtime

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
)

// Nil represents a nil value.
type Nil struct{}

func (n Nil) String() string { return "nil" }

// Bool represents a boolean value.
type Bool bool

// Equals returns true if 'other' is a boolean and has same logical value.
func (b Bool) Equals(other Value) bool {
	val, ok := other.(Bool)
	return ok && (val == b)
}

func (b Bool) String() string { return fmt.Sprintf("%t", b) }

// Float64 represents double precision floating point numbers represented
// using decimal or scientific number formats.
type Float64 float64

// Equals returns true if 'other' is also a float and has same value.
func (f64 Float64) Equals(other Value) bool {
	val, isFloat := other.(Float64)
	return isFloat && (val == f64)
}

func (f64 Float64) String() string { return fmt.Sprintf("%f", f64) }

// Int64 represents integer values represented using decimal, octal, radix
// and hexadecimal formats.
type Int64 int64

// Equals returns true if the other value is also an integer and has same value.
func (i64 Int64) Equals(other Value) bool {
	val, isInt := other.(Int64)
	return isInt && (val == i64)
}

func (i64 Int64) String() string { return fmt.Sprintf("%d", i64) }

// Char represents a character literal.  For example, \a, \b, \1, \∂ etc are
// valid character literals. In addition, special literals like \newline, \space
// etc are supported by the reader.
type Char rune

// Equals returns true if the other value is also a character and has same value.
func (char Char) Equals(other Value) bool {
	val, isChar := other.(Char)
	return isChar && (val == char)
}

func (char Char) String() string { return fmt.Sprintf("\\%c", rune(char)) }

// String represents double-quoted string literals. String Form represents
// the true string value obtained from the reader. Escape sequences are not
// applicable at this level.
type String string

// Equals returns true if 'other' is string and has same value.
func (str String) Equals(other Value) bool {
	otherStr, isStr := other.(String)
	return isStr && (otherStr == str)
}

func (str String) String() string { return fmt.Sprintf("\"%s\"", string(str)) }

// Keyword represents a keyword literal.
type Keyword string

// Equals returns true if the other value is keyword and has same value.
func (kw Keyword) Equals(other Value) bool {
	otherKW, isKeyword := other.(Keyword)
	return isKeyword && (otherKW == kw)
}

func (kw Keyword) String() string { return fmt.Sprintf(":%s", string(kw)) }

// Symbol represents a name given to a value in memory.
type Symbol struct {
	Position
	Value string
}

// Equals returns true if the other value is also a symbol and has same value.
func (sym Symbol) Equals(other Value) bool {
	otherSym, isSym := other.(Symbol)
	return isSym && (sym.Value == otherSym.Value)
}

func (sym Symbol) String() string { return sym.Value }
