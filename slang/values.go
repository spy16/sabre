package slang

import (
	"strings"

	"github.com/spy16/sabre"
)

// MakeString returns stringified version of all args.
func MakeString(vals ...sabre.Value) sabre.Value {
	argc := len(vals)
	switch argc {
	case 0:
		return sabre.String("")

	case 1:
		nilVal := sabre.Nil{}
		if vals[0] == nilVal || vals[0] == nil {
			return sabre.String("")
		}

		return sabre.String(strings.Trim(vals[0].String(), "\""))

	default:
		var sb strings.Builder
		for _, v := range vals {
			sb.WriteString(strings.Trim(v.String(), "\""))
		}
		return sabre.String(sb.String())
	}
}
