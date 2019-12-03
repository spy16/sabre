package sabre

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"unicode"
)

var (
	// ErrSkip is returned by Reader when a no-op form is obtained to indicate
	// it should be discarded.
	ErrSkip = errors.New("skip expr")

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

// New returns a lisp reader instance which can read forms from r. Reader
// behavior can be customized by using SetMacro to override or remove from
// the default read table.
func New(fileName string, rs io.Reader) *Reader {
	rd := &Reader{
		Stream: Stream{
			File: fileName,
			rs:   bufio.NewReader(rs),
		},
		macros: defaultReadTable(),
	}

	return rd
}

// ReaderMacro implementations can be plugged into the Reader to extend, override
// or customize behavior of the reader.
type ReaderMacro func(rd *Reader, init rune) (Form, error)

// Reader provides functions to parse characters from a stream into symbolic
// expressions or forms.
type Reader struct {
	Stream

	macros map[rune]ReaderMacro
}

// All consumes characters from stream until EOF and returns a list of all the
// forms parsed. Any no-op forms (e.g., comment) returned will not be included
// in the result.
func (rd *Reader) All() (Module, error) {
	var forms []Form

	for {
		form, err := rd.One()
		if err != nil {
			if err == ErrSkip {
				continue
			} else if err == io.EOF {
				break
			}

			return nil, err
		}

		forms = append(forms, form)
	}

	return forms, nil
}

// One consumes characters from underlying stream until a complete form is
// parsed and returns the form. In case of no-op forms like comments, this
// returns nil form with ErrSkip. Except EOF and ErrSkip, all errors will
// be wrapped with ReaderError type along with the positional information
// obtained using Info().
func (rd *Reader) One() (Form, error) {
	form, err := rd.readOne()
	if err != nil {
		return nil, rd.annotateErr(err)
	}

	return form, nil
}

// IsTerminal returns true if the rune should terminate a form. ReaderMacro
// trigger runes defined in the read table and all space characters including
// "," are considered terminal.
func (rd *Reader) IsTerminal(r rune) bool {
	_, found := rd.macros[r]
	return found || isSpace(r)
}

// SetMacro sets the given reader macro as the handler for init rune in the
// read table. Overwrites if a macro is already present. If the macro value
// given is nil, entry for the init rune will be removed from the read table.
func (rd *Reader) SetMacro(init rune, macro ReaderMacro) {
	if macro == nil {
		delete(rd.macros, init)
		return
	}

	rd.macros[init] = macro
}

// readOne is same as One() but always returns un-annotated errors.
func (rd *Reader) readOne() (Form, error) {
	if err := rd.SkipSpaces(); err != nil {
		return nil, err
	}

	r, err := rd.NextRune()
	if err != nil {
		return nil, err
	}

	if unicode.IsNumber(r) {
		return readNumber(rd, r)
	} else if r == '+' || r == '-' {
		r2, err := rd.NextRune()
		if err != nil && err != io.EOF {
			return nil, err
		}

		if err != io.EOF {
			rd.Unread(r2)
			if unicode.IsNumber(r2) {
				return readNumber(rd, r)
			}
		}
	}

	macro, found := rd.macros[r]
	if found {
		return macro(rd, r)
	}

	return readSymbol(rd, r)
}

func (rd *Reader) annotateErr(e error) error {
	if e == io.EOF || e == ErrSkip {
		return e
	}

	file, line, col := rd.Info()
	return ReaderError{
		Cause:  e,
		File:   file,
		Line:   line,
		Column: col,
	}
}

func readNumber(rd *Reader, init rune) (Form, error) {
	numStr, err := readToken(rd, init)
	if err != nil {
		return nil, err
	}
	decimalPoint := strings.ContainsRune(numStr, '.')
	isRadix := strings.ContainsRune(numStr, 'r')
	isScientific := strings.ContainsRune(numStr, 'e')

	numErr := fmt.Errorf("illegal number format: '%s'", numStr)

	if isRadix && (decimalPoint || isScientific) {
		return nil, numErr
	}

	var num Number

	if isScientific {
		parts := strings.Split(numStr, "e")
		if len(parts) != 2 {
			return nil, numErr
		}

		base, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return nil, numErr
		}

		pow, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return nil, numErr
		}

		num.IsFloat = true
		num.Float = base * math.Pow(10, float64(pow))
	} else if decimalPoint {
		num.IsFloat = true
		num.Float, err = strconv.ParseFloat(numStr, 64)
	} else {
		base := int64(0)
		repr := numStr

		if isRadix {
			parts := strings.Split(numStr, "r")
			if len(parts) != 2 {
				return nil, numErr
			}

			base, err = strconv.ParseInt(parts[0], 10, 64)
			if err != nil {
				return nil, numErr
			}

			repr = parts[1]
			if base < 0 {
				base = -1 * base
				repr = "-" + repr
			}
		}

		num.Int, err = strconv.ParseInt(repr, int(base), 64)
	}

	if err != nil {
		return nil, numErr
	}

	return num, nil
}

func readSymbol(rd *Reader, init rune) (Form, error) {
	s, err := readToken(rd, init)
	if err != nil {
		return nil, err
	}

	return Symbol(s), nil
}

func readKeyword(rd *Reader, init rune) (Form, error) {
	token, err := readToken(rd, init)
	if err != nil {
		return nil, err
	}

	return Keyword(token), nil
}

func readCharacter(rd *Reader, _ rune) (Form, error) {
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

func readString(rd *Reader, _ rune) (Form, error) {
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

func readList(rd *Reader, _ rune) (Form, error) {
	forms, err := readContainer(rd, '(', ')', "list")
	if err != nil {
		return nil, err
	}

	return List{Forms: forms}, nil
}

func readVector(rd *Reader, _ rune) (Form, error) {
	forms, err := readContainer(rd, '[', ']', "vector")
	if err != nil {
		return nil, err
	}

	return Vector(forms), nil
}

func readComment(rd *Reader, _ rune) (Form, error) {
	for {
		r, err := rd.NextRune()
		if err != nil {
			return nil, err
		}

		if r == '\n' {
			break
		}
	}

	return nil, ErrSkip
}

func quoteFormReader(expandFunc string) ReaderMacro {
	return func(rd *Reader, _ rune) (Form, error) {
		expr, err := rd.One()
		if err != nil {
			if err == io.EOF {
				return nil, errors.New("EOF while reading quote form")
			} else if err == ErrSkip {
				return nil, errors.New("no-op form while reading quote form")
			}
			return nil, err
		}

		return List{
			Forms: []Form{
				Symbol(expandFunc),
				expr,
			},
		}, nil
	}
}

func unmatchedDelimiter(_ *Reader, initRune rune) (Form, error) {
	return nil, fmt.Errorf("unmatched delimiter '%c'", initRune)
}

func readToken(rd *Reader, init rune) (string, error) {
	var b strings.Builder
	b.WriteRune(init)

	for {
		r, err := rd.NextRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		if rd.IsTerminal(r) {
			rd.Unread(r)
			break
		}

		b.WriteRune(r)
	}

	return b.String(), nil
}

func readContainer(rd *Reader, _ rune, end rune, formType string) ([]Form, error) {
	var forms []Form

	for {
		if err := rd.SkipSpaces(); err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("EOF while reading %s", formType)
			}
			return nil, err
		}

		r, err := rd.NextRune()
		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("EOF while reading %s", formType)
			}
			return nil, err
		}

		if r == end {
			break
		}
		rd.Unread(r)

		expr, err := rd.One()
		if err != nil {
			if err == ErrSkip {
				continue
			}
			return nil, err
		}
		forms = append(forms, expr)
	}

	return forms, nil
}

func getEscape(r rune) (rune, error) {
	escaped, found := escapeMap[r]
	if !found {
		return -1, fmt.Errorf("illegal escape sequence '\\%c'", r)
	}

	return escaped, nil
}

func defaultReadTable() map[rune]ReaderMacro {
	return map[rune]ReaderMacro{
		'"':  readString,
		';':  readComment,
		'(':  readList,
		')':  unmatchedDelimiter,
		'[':  readVector,
		']':  unmatchedDelimiter,
		':':  readKeyword,
		'\\': readCharacter,
		'\'': quoteFormReader("quote"),
		'~':  quoteFormReader("unquote"),
	}
}

// ReaderError wraps the parsing error with file and positional information.
type ReaderError struct {
	Cause  error
	File   string
	Line   int
	Column int
}

func (err ReaderError) Error() string {
	return fmt.Sprintf("syntax error in '%s' (Line %d Col %d): %v", err.File, err.Line, err.Column, err.Cause)
}
