package sabre

import (
	"fmt"
	"io"
	"strings"

	"github.com/spy16/sabre/sabre/reader"
	"github.com/spy16/sabre/sabre/runtime"
)

var _ runtime.Runtime = (*Sabre)(nil)

// New returns a new root Sabre instance.
func New() *Sabre {
	rt := runtime.New(nil)
	// TODO: add bindings for builtins
	return &Sabre{Runtime: rt}
}

// ReadEval reads forms from 'r' until EOF using the default reader instance and
// evaluates all forms against the runtime. Returns the result of the last eval.
func ReadEval(rt runtime.Runtime, r io.Reader) (runtime.Value, error) {
	forms, err := reader.New(r).All()
	if err != nil {
		return nil, err
	}

	if len(forms) == 0 {
		return runtime.Nil{}, nil
	}

	res, err := runtime.EvalAll(rt, forms)
	if err != nil {
		return nil, err
	}

	return res[len(res)-1], nil
}

// Sabre implements a sabre runtime with support for qualified symbols.
type Sabre struct {
	runtime.Runtime
}

// Bind creates a binding for symbol to value. Bind throws error if the
// symbol is qualified and the target value does not support attributes.
func (s *Sabre) Bind(symbol string, v runtime.Value) error {
	// TODO: Resolve target and check if it is Attributable.
	if strings.Contains(symbol, ".") {
		return fmt.Errorf("cannot bind to qualified symbol")
	}
	return s.Runtime.Bind(symbol, v)
}

// Resolve resolves the fully-qualified symbol and returns the value.
func (s *Sabre) Resolve(symbol string) (runtime.Value, error) {
	fields := strings.SplitN(symbol, ".", 2)

	if symbol == "." {
		fields = []string{"."}
	}

	target, err := s.Resolve(fields[0])
	if len(fields) == 1 || err != nil {
		return target, err
	}

	return AccessMember(target, fields[1:])
}
