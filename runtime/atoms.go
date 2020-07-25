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

	_ Invokable = Keyword("specimen")
)

// Nil represents a nil value.
type Nil struct{}

// Eval returns the underlying value.
func (n Nil) Eval(_ Runtime) (Value, error) { return n, nil }

func (n Nil) String() string { return "nil" }

// Bool represents a boolean value.
type Bool bool

// Eval returns the underlying value.
func (b Bool) Eval(_ Runtime) (Value, error) { return b, nil }

// Equals returns true if 'other' is a boolean and has same logical value.
func (b Bool) Equals(other Value) bool {
	val, ok := other.(Bool)
	return ok && (val == b)
}

func (b Bool) String() string { return fmt.Sprintf("%t", b) }

// Float64 represents double precision floating point numbers represented
// using decimal or scientific number formats.
type Float64 float64

// Eval simply returns itself since Floats evaluate to themselves.
func (f64 Float64) Eval(_ Runtime) (Value, error) { return f64, nil }

// Equals returns true if 'other' is also a float and has same value.
func (f64 Float64) Equals(other Value) bool {
	val, isFloat := other.(Float64)
	return isFloat && (val == f64)
}

func (f64 Float64) String() string { return fmt.Sprintf("%f", f64) }

// Int64 represents integer values represented using decimal, octal, radix
// and hexadecimal formats.
type Int64 int64

// Eval simply returns itself since Integers evaluate to themselves.
func (i64 Int64) Eval(_ Runtime) (Value, error) { return i64, nil }

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

// Eval simply returns itself since Chracters evaluate to themselves.
func (char Char) Eval(_ Runtime) (Value, error) { return char, nil }

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

// Eval simply returns itself since Strings evaluate to themselves.
func (se String) Eval(_ Runtime) (Value, error) { return se, nil }

// Equals returns true if 'other' is string and has same value.
func (se String) Equals(other Value) bool {
	val, isStr := other.(String)
	return isStr && (val == se)
}

func (se String) String() string { return fmt.Sprintf("\"%s\"", string(se)) }

// Keyword represents a keyword literal.
type Keyword string

// Eval simply returns itself since Keywords evaluate to themselves.
func (kw Keyword) Eval(_ Runtime) (Value, error) { return kw, nil }

// Equals returns true if the other value is keyword and has same value.
func (kw Keyword) Equals(other Value) bool {
	val, isKeyword := other.(Keyword)
	return isKeyword && (val == kw)
}

func (kw Keyword) String() string { return fmt.Sprintf(":%s", string(kw)) }

// Invoke enables keyword lookup for maps.
func (kw Keyword) Invoke(scope Runtime, args ...Value) (Value, error) {
	if len(args) != 1 && len(args) != 2 {
		return nil, fmt.Errorf("keyword specialInvoke requires 1 or 2 arguments, got %d", len(args))
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

	val := assocVal.EntryAt(kw)
	if val == nil {
		val = def
	}

	return val, nil
}

// Symbol represents a name given to a value in memory.
type Symbol struct {
	Position
	Value string
}

// Eval returns the value bound to this symbol in current context.
func (sym Symbol) Eval(scope Runtime) (Value, error) {
	return scope.Resolve(sym.Value)
}

// Equals returns true if the other value is also a symbol and has same value.
func (sym Symbol) Equals(other Value) bool {
	val, isSym := other.(Symbol)
	return isSym && (sym.Value == val.Value)
}

func (sym Symbol) String() string { return sym.Value }
