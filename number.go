package sabre

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

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

func readNumber(rd *Reader, init rune) (Value, error) {
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
		return Float64(v), nil

	case isRadix:
		return parseRadix(numStr)

	default:
		v, err := strconv.ParseInt(numStr, 0, 64)
		if err != nil {
			return nil, fmt.Errorf("illegal number format '%s'", numStr)
		}

		return Int64(v), nil
	}
}

func parseRadix(numStr string) (Int64, error) {
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

	return Int64(v), nil
}

func parseScientific(numStr string) (Float64, error) {
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

	return Float64(base * math.Pow(10, float64(pow))), nil
}
