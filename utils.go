package sabre

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"reflect"
	"strings"
	"unicode"
)

func defaultReadTable() map[rune]Macro {
	return map[rune]Macro{
		'"':  readString,
		';':  readComment,
		':':  readKeyword,
		'\\': readCharacter,
		'\'': quoteFormReader("quote"),
		'~':  quoteFormReader("unquote"),
		'(':  readList,
		')':  unmatchedDelimiter,
		'[':  readVector,
		']':  unmatchedDelimiter,
	}
}

func containerString(vals []Value, begin, end, sep string) string {
	parts := make([]string, len(vals))
	for i, expr := range vals {
		parts[i] = fmt.Sprintf("%v", expr)
	}
	return begin + strings.Join(parts, sep) + end
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
