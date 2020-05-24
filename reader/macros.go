package reader

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/spy16/sabre/core"
)

// Macro implementations can be plugged into the Reader to extend, override
// or customize behavior of the reader.
type Macro func(rd *Reader, init rune) (core.Value, error)

func defaultReadTable() map[rune]Macro {
	return map[rune]Macro{
		'"':  readString,
		';':  readComment,
		':':  readKeyword,
		'\\': readCharacter,
		'\'': quoteFormReader("quote"),
		'~':  quoteFormReader("unquote"),
		'`':  quoteFormReader("syntax-quote"),
		'(':  readList,
		')':  unmatchedDelimiter,
		'[':  readVector,
		']':  unmatchedDelimiter,
		'{':  readHashMap,
		'}':  unmatchedDelimiter,
	}
}

func defaultDispatchTable() map[rune]Macro {
	return map[rune]Macro{
		'{': readSet,
		'}': unmatchedDelimiter,
	}
}

func readString(rd *Reader, _ rune) (core.Value, error) {
	var b strings.Builder

	for {
		r, err := rd.NextRune()
		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("%w: while reading string", ErrEOF)
			}

			return nil, err
		}

		if r == '\\' {
			r2, err := rd.NextRune()
			if err != nil {
				if err == io.EOF {
					return nil, fmt.Errorf("%w: while reading string", ErrEOF)
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

	return core.String(b.String()), nil
}

func readNumber(rd *Reader, init rune) (core.Value, error) {
	numStr, err := readToken(rd, init)
	if err != nil {
		return nil, err
	}

	decimalPoint := strings.ContainsRune(numStr, '.')
	isRadix := strings.ContainsRune(numStr, 'r')
	isScientific := strings.ContainsRune(numStr, 'e')

	switch {
	case isRadix && (decimalPoint || isScientific):
		return nil, fmt.Errorf("illegal number format: '%s'", numStr)

	case isScientific:
		return parseScientific(numStr)

	case decimalPoint:
		v, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return nil, fmt.Errorf("illegal number format: '%s'", numStr)
		}
		return core.Float64(v), nil

	case isRadix:
		return parseRadix(numStr)

	default:
		v, err := strconv.ParseInt(numStr, 0, 64)
		if err != nil {
			return nil, fmt.Errorf("illegal number format '%s'", numStr)
		}

		return core.Int64(v), nil
	}
}

func readSymbol(rd *Reader, init rune) (core.Value, error) {
	pi := rd.Position()

	s, err := readToken(rd, init)
	if err != nil {
		return nil, err
	}

	return core.Symbol{
		Value:    s,
		Position: pi,
	}, nil
}

func readKeyword(rd *Reader, init rune) (core.Value, error) {
	token, err := readToken(rd, -1)
	if err != nil {
		return nil, err
	}

	return core.Keyword(token), nil
}

func readCharacter(rd *Reader, _ rune) (core.Value, error) {
	r, err := rd.NextRune()
	if err != nil {
		return nil, fmt.Errorf("%w: while reading character", ErrEOF)
	}

	token, err := readToken(rd, r)
	if err != nil {
		return nil, err
	}
	runes := []rune(token)

	if len(runes) == 1 {
		return core.Character(runes[0]), nil
	}

	v, found := charLiterals[token]
	if found {
		return core.Character(v), nil
	}

	if token[0] == 'u' {
		return readUnicodeChar(token[1:], 16)
	}

	return nil, fmt.Errorf("unsupported character: '\\%s'", token)
}

func readList(rd *Reader, _ rune) (core.Value, error) {
	pi := rd.Position()
	forms, err := readContainer(rd, '(', ')', "list")
	if err != nil {
		return nil, err
	}

	return &core.List{
		Values:   forms,
		Position: pi,
	}, nil
}

func readHashMap(rd *Reader, _ rune) (core.Value, error) {
	pi := rd.Position()
	forms, err := readContainer(rd, '{', '}', "hash-map")
	if err != nil {
		return nil, err
	}

	if len(forms)%2 != 0 {
		return nil, errors.New("expecting even number of forms within {}")
	}

	hm := &core.HashMap{
		Position: pi,
		Data:     map[core.Value]core.Value{},
	}

	for i := 0; i < len(forms); i += 2 {
		if !core.IsHashable(forms[i]) {
			return nil, fmt.Errorf("value of type '%s' is not hashable",
				reflect.TypeOf(forms[i]))
		}

		hm.Data[forms[i]] = forms[i+1]
	}

	return hm, nil
}

func readVector(rd *Reader, _ rune) (core.Value, error) {
	pi := rd.Position()

	forms, err := readContainer(rd, '[', ']', "vector")
	if err != nil {
		return nil, err
	}

	return core.Vector{
		Values:   forms,
		Position: pi,
	}, nil
}

func readSet(rd *Reader, _ rune) (core.Value, error) {
	pi := rd.Position()

	forms, err := readContainer(rd, '{', '}', "set")
	if err != nil {
		return nil, err
	}

	set := core.Set{
		Values:   forms,
		Position: pi,
	}
	if !set.Valid() {
		return nil, errors.New("duplicate value in set")
	}

	return set, nil
}

func readUnicodeChar(token string, base int) (core.Character, error) {
	num, err := strconv.ParseInt(token, base, 64)
	if err != nil {
		return -1, fmt.Errorf("invalid unicode character: '\\%s'", token)
	}

	if num < 0 || num >= unicode.MaxRune {
		return -1, fmt.Errorf("invalid unicode character: '\\%s'", token)
	}

	return core.Character(num), nil
}

func readComment(rd *Reader, _ rune) (core.Value, error) {
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

func unmatchedDelimiter(_ *Reader, initRune rune) (core.Value, error) {
	return nil, fmt.Errorf("unmatched delimiter '%c'", initRune)
}
