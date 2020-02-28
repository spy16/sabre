package sabre

import (
	"fmt"
	"reflect"
	"strings"
)

var nilValue = Nil{}

// Nil represents a nil value.
type Nil struct{}

// Eval returns the underlying value.
func (n Nil) Eval(_ Scope) (Value, error) { return n, nil }

func (n Nil) String() string { return "nil" }

// Bool represents a boolean value.
type Bool bool

// Eval returns the underlying value.
func (b Bool) Eval(_ Scope) (Value, error) { return b, nil }

func (b Bool) String() string { return fmt.Sprintf("%t", b) }

// Float64 represents double precision floating point numbers represented
// using decimal or scientific number formats.
type Float64 float64

// Eval simply returns itself since Floats evaluate to themselves.
func (f64 Float64) Eval(_ Scope) (Value, error) { return f64, nil }

func (f64 Float64) String() string { return fmt.Sprintf("%f", f64) }

// Int64 represents integer values represented using decimal, octal, radix
// and hexadecimal formats.
type Int64 int64

// Eval simply returns itself since Integers evaluate to themselves.
func (i64 Int64) Eval(_ Scope) (Value, error) { return i64, nil }

func (i64 Int64) String() string { return fmt.Sprintf("%d", i64) }

// String represents double-quoted string literals. String Form represents
// the true string value obtained from the reader. Escape sequences are not
// applicable at this level.
type String string

// Eval simply returns itself since Strings evaluate to themselves.
func (se String) Eval(_ Scope) (Value, error) { return se, nil }

func (se String) String() string { return fmt.Sprintf("\"%s\"", string(se)) }

// First returns the first character if string is not empty, nil otherwise.
func (se String) First() Value {
	if len(se) == 0 {
		return nilValue
	}

	return Character(se[0])
}

// Next slices the string by excluding first character and returns the
// remainder.
func (se String) Next() Seq { return se.chars().Next() }

// Cons converts the string to character sequence and adds the given value
// to the beginning of the list.
func (se String) Cons(v Value) Seq { return se.chars().Cons(v) }

// Conj joins the given values to list of characters of the string and returns
// the new sequence.
func (se String) Conj(vals ...Value) Seq { return se.chars().Conj(vals...) }

func (se String) chars() Values {
	var vals Values
	for _, r := range se {
		vals = append(vals, Character(r))
	}
	return vals
}

// Character represents a character literal.  For example, \a, \b, \1, \∂ etc
// are valid character literals. In addition, special literals like \newline,
// \space etc are supported by the reader.
type Character rune

// Eval simply returns itself since Chracters evaluate to themselves.
func (char Character) Eval(_ Scope) (Value, error) { return char, nil }

func (char Character) String() string { return fmt.Sprintf("\\%c", rune(char)) }

// Keyword represents a keyword literal.
type Keyword string

// Eval simply returns itself since Keywords evaluate to themselves.
func (kw Keyword) Eval(_ Scope) (Value, error) { return kw, nil }

func (kw Keyword) String() string { return fmt.Sprintf(":%s", string(kw)) }

// Symbol represents a name given to a value in memory.
type Symbol struct {
	Position
	Value string
}

// Eval returns the value bound to this symbol in current context. If the
// symbol is in fully qualified form (i.e., separated by '.'), eval does
// recursive member access.
func (sym Symbol) Eval(scope Scope) (Value, error) {
	fields := strings.Split(sym.Value, ".")

	if sym.Value == "." {
		fields = []string{"."}
	}

	target, err := scope.Resolve(fields[0])
	if len(fields) == 1 || err != nil {
		return target, err
	}

	rv := reflect.ValueOf(target)
	for i := 1; i < len(fields); i++ {
		if rv.Type() == reflect.TypeOf(Any{}) {
			rv = rv.Interface().(Any).V
		}

		rv, err = accessMember(rv, fields[i])
		if err != nil {
			return nil, err
		}
	}

	if isKind(rv.Type(), reflect.Chan, reflect.Array,
		reflect.Func, reflect.Ptr) && rv.IsNil() {
		return Nil{}, nil
	}

	return ValueOf(rv.Interface()), nil
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
