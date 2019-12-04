package sabre

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

var (
	errStringEOF = errors.New("EOF while reading string")
	errCharEOF   = errors.New("EOF while reading character")
)

var (
	escapeMap = map[rune]rune{
		'"':  '"',
		'n':  '\n',
		'\\': '\\',
		't':  '\t',
		'a':  '\a',
		'f':  '\a',
		'r':  '\r',
		'b':  '\b',
		'v':  '\v',
	}

	charLiterals = map[string]rune{
		"tab":       '\t',
		"space":     ' ',
		"newline":   '\n',
		"return":    '\r',
		"backspace": '\b',
		"formfeed":  '\f',
	}
)

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

func readString(rd *Reader, _ rune) (Value, error) {
	var b strings.Builder

	for {
		r, err := rd.NextRune()
		if err != nil {
			if err == io.EOF {
				return nil, errStringEOF
			}

			return nil, err
		}

		if r == '\\' {
			r2, err := rd.NextRune()
			if err != nil {
				if err == io.EOF {
					return nil, errStringEOF
				}

				return nil, err
			}

			// TODO: Support for Unicode escape \uNN format.

			escaped, err := getEscape(r2)
			if err != nil {
				return nil, err
			}
			r = escaped
		} else if r == '"' {
			break
		}

		b.WriteRune(r)
	}

	return String(b.String()), nil
}

func readSymbol(rd *Reader, init rune) (Value, error) {
	s, err := readToken(rd, init)
	if err != nil {
		return nil, err
	}

	return Symbol(s), nil
}

func readKeyword(rd *Reader, init rune) (Value, error) {
	token, err := readToken(rd, init)
	if err != nil {
		return nil, err
	}

	return Keyword(token), nil
}

func readCharacter(rd *Reader, _ rune) (Value, error) {
	r, err := rd.NextRune()
	if err != nil {
		return nil, errCharEOF
	}

	token, err := readToken(rd, r)
	if err != nil {
		return nil, err
	}
	runes := []rune(token)

	if len(runes) == 1 {
		return Character(runes[0]), nil
	}

	v, found := charLiterals[token]
	if found {
		return Character(v), nil
	}

	if token[0] == 'u' {
		return readUnicodeChar(token[1:], 16)
	}

	return nil, fmt.Errorf("unsupported character: '\\%s'", token)
}

func readUnicodeChar(token string, base int) (Character, error) {
	num, err := strconv.ParseInt(token, base, 64)
	if err != nil {
		return -1, fmt.Errorf("invalid unicode character: '\\%s'", token)
	}

	if num < 0 || num >= unicode.MaxRune {
		return -1, fmt.Errorf("invalid unicode character: '\\%s'", token)
	}

	return Character(num), nil
}

func getEscape(r rune) (rune, error) {
	escaped, found := escapeMap[r]
	if !found {
		return -1, fmt.Errorf("illegal escape sequence '\\%c'", r)
	}

	return escaped, nil
}
