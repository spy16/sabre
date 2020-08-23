package runtime

import (
	"fmt"
	"reflect"
	"strings"
)

// Invoke evaluates the first item in the seq and invokes the resultant
// value with rest of the list as arguments.
func Invoke(rt Runtime, seq Seq) (Value, error) {
	if seq.Count() == 0 {
		return seq, nil
	}

	v, err := rt.Eval(seq.First())
	if err != nil {
		return nil, err
	}

	target, ok := v.(Invokable)
	if !ok {
		return nil, fmt.Errorf("value of type '%s' is not invokable", reflect.TypeOf(v))
	}

	return target.Invoke(rt, toSlice(seq.Next())...)
}

// EvalAll evaluates each value in the list against the given env and returns a list
// of resultant value.
func EvalAll(rt Runtime, vals []Value) ([]Value, error) {
	var results []Value
	for _, f := range vals {
		res, err := rt.Eval(f)
		if err != nil {
			return nil, err
		}
		results = append(results, res)
	}
	return results, nil
}

// Equals compares two values in an identity independent manner. If v1 implements
// `Equals(Value)` method, then the comparison is delegated to it.
func Equals(v1, v2 Value) bool {
	if isNil(v1) && isNil(v2) {
		return true
	}

	if cmp, ok := v1.(interface{ Equals(other Value) bool }); ok {
		return cmp.Equals(v2)
	}

	s1, isV1Seq := v1.(Seq)
	s2, isV2Seq := v2.(Seq)
	if isV1Seq && isV2Seq {
		return compareSeq(s1, s2)
	}

	return reflect.DeepEqual(v1, v2)
}

// Cons returns a new seq with `v` added as the first and `seq` as the rest. Seq
// can be nil as well.
func Cons(v Value, seq Seq) Seq {
	newSeq := &LinkedList{
		first: v,
		rest:  seq,
		count: 1,
	}

	if seq != nil {
		newSeq.count = seq.Count() + 1
	}

	return newSeq
}

// SeqString returns a string representation for the sequence with given prefix
// suffix and separator.
func SeqString(seq Seq, begin, end, sep string) string {
	var parts []string
	ForEach(seq, func(item Value) bool {
		parts = append(parts, item.String())
		return false
	})
	return begin + strings.Join(parts, sep) + end
}

// ForEach reads from the sequence and calls the given function for each item.
// Function can return true to stop the iteration.
func ForEach(seq Seq, call func(item Value) bool) {
	for seq != nil {
		v := seq.First()
		if v == nil || call(seq.First()) {
			break
		}
		seq = seq.Next()
	}
}

func toSlice(seq Seq) []Value {
	var slice []Value
	ForEach(seq, func(item Value) bool {
		slice = append(slice, item)
		return false
	})
	return slice
}

func isNil(v Value) bool {
	_, isNil := v.(Nil)
	return v == nil || isNil
}

func compareSeq(s1, s2 Seq) bool {
	if s1.Count() != s2.Count() {
		return false
	}
	for s1 != nil && s2 != nil {
		if !Equals(s1.First(), s2.First()) {
			return false
		}
		s1 = s1.Next()
		s2 = s2.Next()
	}
	return true
}

// Position represents the positional information about a value read
// by reader.
type Position struct {
	File   string
	Line   int
	Column int
}

// Pos returns the file, line and column values.
func (pi Position) Pos() (file string, line, col int) {
	return pi.File, pi.Line, pi.Column
}

// SetPos sets the position information.
func (pi *Position) SetPos(file string, line, col int) {
	pi.File = file
	pi.Line = line
	pi.Column = col
}

func (pi Position) String() string {
	if pi.File == "" {
		pi.File = "<unknown>"
	}
	return fmt.Sprintf("%s:%d:%d", pi.File, pi.Line, pi.Column)
}
