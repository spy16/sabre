package sabre

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spy16/sabre/sabre/core"
)

// Sabre implements an instance of sabre interpreter.
type Sabre struct {
}

func (s *Sabre) Eval(form core.Value) (core.Value, error) {
	panic("not implemented") // TODO: Implement
}

func (s *Sabre) Bind(symbol string, v core.Value) error {
	panic("not implemented") // TODO: Implement
}

func (s *Sabre) Resolve(symbol string) (core.Value, error) {
	fields := strings.SplitN(symbol, ".", 2)

	if symbol == "." {
		fields = []string{"."}
	}

	target, err := s.Resolve(fields[0])
	if len(fields) == 1 || err != nil {
		return target, err
	}

	attr, ok := target.(core.Attributable)
	if !ok {
		return nil, fmt.Errorf("attributes not supported on '%s'", reflect.TypeOf(target))
	}

	val := attr.GetAttr(fields[1], nil)
	if val == nil {
		return nil, fmt.Errorf(
			"value of type '%s' does not have attribute '%s'",
			reflect.TypeOf(attr), fields[1],
		)
	}

	return val, nil
}

func (s *Sabre) Parent() core.Env {
	panic("not implemented") // TODO: Implement
}
