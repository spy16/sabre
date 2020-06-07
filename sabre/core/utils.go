package core

import (
	"strings"

	"github.com/spy16/sabre/sabre/runtime"
)

// Module is a convenience wrapper for evaluating multiple values. Evaluation
// returns last result of evaluating all forms in the module.
type Module []runtime.Value

// Eval evaluates each form in the module and returns the result of last eval.
// Returns Nil if the module is empty.
func (mod Module) Eval(rt runtime.Runtime) (runtime.Value, error) {
	var res runtime.Value
	var err error
	for _, f := range mod {
		res, err = rt.Eval(f)
		if err != nil {
			return nil, err
		}
	}

	if runtime.Equals(res, nil) {
		return runtime.Nil{}, nil
	}

	return res, nil
}

// String returns module represented using a (do expr*) form.
func (mod Module) String() string {
	var sb strings.Builder
	sb.WriteString("(do ")
	for i, f := range mod {
		sb.WriteString(f.String())
		if i < len(mod)-1 {
			sb.WriteString("\n    ")
		}
	}
	sb.WriteRune(')')
	return sb.String()
}
