package sabre

import (
	"fmt"
	"io"
	"strings"

	"github.com/spy16/sabre/core"
	"github.com/spy16/sabre/reader"
	"github.com/spy16/sabre/runtime"
)

const sep = "."

var _ runtime.Runtime = (*Sabre)(nil)

// ValueOf converts arbitrary Go value to runtime value. This function is an alias
// for core.ValueOf() provided for convenience.
var ValueOf = core.ValueOf

// New returns a new root Sabre instance with built-in special forms.
func New() *Sabre {
	rt := runtime.New(nil)
	_ = rt.Bind("do", runtime.GoFunc(doForm))
	_ = rt.Bind("let", runtime.GoFunc(letForm))
	_ = rt.Bind("def", runtime.GoFunc(defForm))
	_ = rt.Bind("cond", runtime.GoFunc(condForm))
	return &Sabre{Runtime: rt}
}

// ReadEval reads forms from 'r' until EOF using the default reader instance and
// evaluates all forms against the runtime. Returns the result of the last eval.
func ReadEval(rt runtime.Runtime, r io.Reader) (runtime.Value, error) {
	forms, err := reader.New(r).All()
	if err != nil {
		return nil, err
	}
	return rt.Eval(core.Module(forms))
}

// Sabre implements a sabre runtime with support for qualified symbols.
type Sabre struct {
	runtime.Runtime
}

// Bind creates a binding for symbol to value. Bind throws error if the symbol is
// qualified and the target value does not support attributes.
func (s *Sabre) Bind(symbol string, v runtime.Value) error {
	// TODO: Resolve target and check if it is Attributable.
	if strings.Contains(symbol, sep) {
		return fmt.Errorf("cannot bind to qualified symbol")
	}
	return s.Runtime.Bind(symbol, v)
}

// Resolve recursively resolves the fully-qualified symbol and returns the value.
func (s *Sabre) Resolve(symbol string) (runtime.Value, error) {
	fields := strings.SplitN(symbol, sep, 2)

	if symbol == sep {
		fields = []string{sep}
	}

	target, err := s.Runtime.Resolve(fields[0])
	if len(fields) == 1 || err != nil {
		return target, err
	}

	return core.AccessMember(target, fields[1:])
}
