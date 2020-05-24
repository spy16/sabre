package reader

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/spy16/sabre/core"
)

const dispatchTrigger = '#'

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

	predefSymbols = map[string]core.Value{
		"nil":   core.Nil{},
		"true":  core.Bool(true),
		"false": core.Bool(false),
	}
)

// New returns a lisp reader instance which can read forms from rs. Reader
// behavior can be customized by using SetMacro to override or remove from
// the default read table. File name will be inferred from the reader value
// and type information or can be set manually on the Reader.
func New(rs io.Reader) *Reader {
	return &Reader{
		File:     inferFileName(rs),
		rs:       bufio.NewReader(rs),
		macros:   defaultReadTable(),
		dispatch: defaultDispatchTable(),
	}
}

// Reader provides functions to parse characters from a stream into symbolic
// expressions or forms.
type Reader struct {
	File string

	rs          io.RuneReader
	buf         []rune
	line, col   int
	lastCol     int
	macros      map[rune]Macro
	dispatch    map[rune]Macro
	dispatching bool
}

// All consumes characters from stream until EOF and returns a list of all the
// forms parsed. Any no-op forms (e.g., comment) returned will not be included
// in the result.
func (rd *Reader) All() (core.Value, error) {
	var forms []core.Value

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

	return core.Module(forms), nil
}

// One consumes characters from  underlying stream until a  complete form is
// parsed and returns  the form while ignoring the no-op forms like comments.
// Except  EOF, all errors will  be wrapped with ReaderError type along with
// the positional information obtained using Position().
func (rd *Reader) One() (core.Value, error) {
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

// IsTerminal returns  true if the rune  should terminate a form. ReaderMacro
// trigger runes defined in the read table and all space characters including
// "," are considered terminal.
func (rd *Reader) IsTerminal(r rune) bool {
	if isSpace(r) {
		return true
	}

	if rd.dispatching {
		_, found := rd.dispatch[r]
		if found {
			return true
		}
	}

	_, found := rd.macros[r]
	return found
}

// SetMacro sets the given reader  macro as the handler for init rune in the
// read table. Overwrites  if a macro is already present. If the macro value
// given is nil, entry for the init rune will be removed from the read table.
// isDispatch decides if the macro is a dispatch macro and takes effect only
// after a '#' sign.
func (rd *Reader) SetMacro(init rune, macro Macro, isDispatch bool) {
	if isDispatch {
		if macro == nil {
			delete(rd.dispatch, init)
			return
		}
		rd.dispatch[init] = macro
	} else {
		if macro == nil {
			delete(rd.macros, init)
			return
		}
		rd.macros[init] = macro
	}
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

// Position returns information about the stream including file name and
// the position of the reader.
func (rd Reader) Position() core.Position {
	file := strings.TrimSpace(rd.File)
	return core.Position{
		File:   file,
		Line:   rd.line + 1,
		Column: rd.col,
	}
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
func (rd *Reader) readOne() (core.Value, error) {
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

	if r == dispatchTrigger {
		f, err := rd.execDispatch()
		if f != nil || err != nil {
			return f, err
		}
	}

	v, err := readSymbol(rd, r)
	if err != nil {
		return nil, err
	}

	if predefVal, found := predefSymbols[v.(core.Symbol).Value]; found {
		return predefVal, nil
	}

	return v, nil
}

func (rd *Reader) execDispatch() (core.Value, error) {
	pos := rd.Position()

	r2, err := rd.NextRune()
	if err != nil {
		// ignore the error and let readOne handle it.
		return nil, nil
	}

	dispatchMacro, found := rd.dispatch[r2]
	if !found {
		rd.Unread(r2)
		return nil, nil
	}

	rd.dispatching = true
	defer func() {
		rd.dispatching = false
	}()

	form, err := dispatchMacro(rd, r2)
	if err != nil {
		return nil, err
	}

	setPosition(form, pos)
	return form, nil
}

func (rd *Reader) annotateErr(e error) error {
	if e == io.EOF || e == ErrSkip {
		return e
	}

	return Error{
		Cause:    e,
		Position: rd.Position(),
	}
}

func quoteFormReader(expandFunc string) Macro {
	return func(rd *Reader, _ rune) (core.Value, error) {
		expr, err := rd.One()
		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("%w: while reading quote form", ErrEOF)
			} else if err == ErrSkip {
				return nil, errors.New("no-op form while reading quote form")
			}
			return nil, err
		}

		return &core.List{
			Values: []core.Value{
				core.Symbol{Value: expandFunc},
				expr,
			},
		}, nil
	}
}

func parseRadix(numStr string) (core.Int64, error) {
	parts := strings.Split(numStr, "r")
	if len(parts) != 2 {
		return 0, fmt.Errorf("illegal radix notation '%s'", numStr)
	}

	base, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("illegal radix notation '%s'", numStr)
	}

	repr := parts[1]
	if base < 0 {
		base = -1 * base
		repr = "-" + repr
	}

	v, err := strconv.ParseInt(repr, int(base), 64)
	if err != nil {
		return 0, fmt.Errorf("illegal radix notation '%s'", numStr)
	}

	return core.Int64(v), nil
}

func parseScientific(numStr string) (core.Float64, error) {
	parts := strings.Split(numStr, "e")
	if len(parts) != 2 {
		return 0, fmt.Errorf("illegal scientific notation '%s'", numStr)
	}

	base, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, fmt.Errorf("illegal scientific notation '%s'", numStr)
	}

	pow, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("illegal scientific notation '%s'", numStr)
	}

	return core.Float64(base * math.Pow(10, float64(pow))), nil
}

func getEscape(r rune) (rune, error) {
	escaped, found := escapeMap[r]
	if !found {
		return -1, fmt.Errorf("illegal escape sequence '\\%c'", r)
	}

	return escaped, nil
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

func readContainer(rd *Reader, _ rune, end rune, formType string) ([]core.Value, error) {
	var forms []core.Value

	for {
		if err := rd.SkipSpaces(); err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("%w: while reading %s", ErrEOF, formType)
			}
			return nil, err
		}

		r, err := rd.NextRune()
		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("%w: while reading %s", ErrEOF, formType)
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

func isSpace(r rune) bool {
	return unicode.IsSpace(r) || r == ','
}

func inferFileName(rs io.Reader) string {
	switch r := rs.(type) {
	case *os.File:
		return r.Name()

	case *strings.Reader:
		return "<string>"

	case *bytes.Reader:
		return "<bytes>"

	case net.Conn:
		return fmt.Sprintf("<con:%s>", r.LocalAddr())

	default:
		return fmt.Sprintf("<%s>", reflect.TypeOf(rs))
	}
}

func setPosition(form core.Value, pos core.Position) core.Value {
	p, canSet := form.(interface {
		SetPos(file string, line, col int)
	})
	if !canSet {
		return form
	}

	p.SetPos(pos.File, pos.Line, pos.Column)
	return form
}
