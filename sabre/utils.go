package sabre

import "strings"

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

func isNil(v Value) bool {
	_, isNil := v.(Nil)
	return v == nil || isNil
}
