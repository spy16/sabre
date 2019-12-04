package sabre

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

// ErrSkip can be returned by reader macro to indicate a no-op form which
// should be discarded (e.g., Comments).
var ErrSkip = errors.New("skip expr")

// NewReader returns a lisp reader instance which can read forms from rs.
// Reader behavior can be customized by using SetMacro to override or remove
// from the default read table. File name  will be inferred from the  reader
// value and type information.
func NewReader(rs io.Reader) *Reader {
	return &Reader{
		File:   inferFileName(rs),
		rs:     bufio.NewReader(rs),
		macros: defaultReadTable(),
	}
}

// Macro implementations can be plugged into the Reader to extend, override
// or customize behavior of the reader.
type Macro func(rd *Reader, init rune) (Value, error)

// Reader provides functions to parse characters from a stream into symbolic
// expressions or forms.
type Reader struct {
	File string
	Hook Macro

	rs        io.RuneReader
	buf       []rune
	line, col int
	lastCol   int
	macros    map[rune]Macro
}

// All consumes characters from stream until EOF and returns a list of all the
// forms parsed. Any no-op forms (e.g., comment) returned will not be included
// in the result.
func (rd *Reader) All() (Value, error) {
	var forms []Value

	for {
		form, err := rd.One()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		forms = append(forms, form)
	}

	return Module(forms), nil
}

// One consumes characters from underlying stream until a complete form is
// parsed and returns the form while ignoring the no-op forms like comments.
// Except EOF, all errors will be wrapped with ReaderError type along with
// the positional information obtained using Info().
func (rd *Reader) One() (Value, error) {
	for {
		form, err := rd.readOne()
		if err != nil {
			if err == ErrSkip {
				continue
			}

			return nil, rd.annotateErr(err)
		}

		return form, nil
	}
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
func (rd *Reader) SetMacro(init rune, macro Macro) {
	if macro == nil {
		delete(rd.macros, init)
		return
	}

	rd.macros[init] = macro
}

// NextRune returns next rune from the stream and advances the stream.
func (rd *Reader) NextRune() (rune, error) {
	var r rune
	if len(rd.buf) > 0 {
		r = rd.buf[0]
		rd.buf = rd.buf[1:]
	} else {
		temp, _, err := rd.rs.ReadRune()
		if err != nil {
			return -1, err
		}

		r = temp
	}

	if r == '\n' {
		rd.line++
		rd.lastCol = rd.col
		rd.col = 0
	} else {
		rd.col++
	}

	return r, nil
}

// Unread can be used to return runes consumed from the stream back to the
// stream. Un-reading more runes than read is guaranteed to work but might
// cause inconsistency in stream positional information.
func (rd *Reader) Unread(runes ...rune) {
	newLine := false
	for _, r := range runes {
		if r == '\n' {
			newLine = true
			break
		}
	}

	if newLine {
		rd.line--
		rd.col = rd.lastCol
	} else {
		rd.col--
	}

	rd.buf = append(runes, rd.buf...)
}

// Info returns information about the stream including file name and the
// position of the reader.
func (rd Reader) Info() (file string, line, col int) {
	file = strings.TrimSpace(rd.File)
	return file, rd.line + 1, rd.col
}

// SkipSpaces consumes and discards runes from stream repeatedly until a
// character that is not a whitespace is identified. Along with standard
// unicode  white-space characters "," is also considered  a white-space
// and discarded.
func (rd *Reader) SkipSpaces() error {
	for {
		r, err := rd.NextRune()
		if err != nil {
			return err
		}

		if !isSpace(r) {
			rd.Unread(r)
			break
		}
	}

	return nil
}

// readOne is same as One() but always returns un-annotated errors.
func (rd *Reader) readOne() (Value, error) {
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

	if rd.Hook != nil {
		f, err := rd.Hook(rd, r)
		if err != ErrSkip {
			return f, err
		}
	}

	return readSymbol(rd, r)
}

func (rd *Reader) annotateErr(e error) error {
	if e == io.EOF || e == ErrSkip {
		return e
	}

	file, line, col := rd.Info()
	return Error{
		Cause:  e,
		File:   file,
		Line:   line,
		Column: col,
	}
}

func readComment(rd *Reader, _ rune) (Value, error) {
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

func quoteFormReader(expandFunc string) Macro {
	return func(rd *Reader, _ rune) (Value, error) {
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
			Symbol(expandFunc),
			expr,
		}, nil
	}
}

func unmatchedDelimiter(_ *Reader, initRune rune) (Value, error) {
	return nil, fmt.Errorf("unmatched delimiter '%c'", initRune)
}

func readToken(rd *Reader, init rune) (string, error) {
	var b strings.Builder
	if init != -1 {
		b.WriteRune(init)
	}

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

func readContainer(rd *Reader, _ rune, end rune, formType string) ([]Value, error) {
	var forms []Value

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

		expr, err := rd.readOne()
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

// Error wraps the parsing error with file and positional information.
type Error struct {
	Cause  error
	File   string
	Line   int
	Column int
}

func (err Error) Error() string {
	if e, ok := err.Cause.(Error); ok {
		return e.Error()
	}

	return fmt.Sprintf("syntax error in '%s' (Line %d Col %d): %v", err.File, err.Line, err.Column, err.Cause)
}
