package core

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

// Env implementations provide the environment for maintaining bindings
// and evaluating forms.
type Env interface {
	Eval(form Value) (Value, error)
	Bind(symbol string, v Value) error
	Resolve(symbol string) (Value, error)
	Parent() Env
}

// New returns an implementation of env with no bindings.
func New(parent Env) Env {
	return &mapEnv{
		mu:     &sync.RWMutex{},
		scope:  map[string]Value{},
		parent: parent,
	}
}

type mapEnv struct {
	mu     *sync.RWMutex
	scope  map[string]Value
	parent Env
}

func (s *mapEnv) Eval(form Value) (Value, error) {
	if form == nil {
		return Nil{}, nil
	}

	if e, ok := form.(Expr); ok {
		res, err := e.Eval(s)
		if err != nil {
			return nil, NewEvalErr(e, err)
		}
		return res, nil
	}

	return form, nil
}

func (s *mapEnv) Bind(symbol string, v Value) error {
	symbol = strings.TrimSpace(symbol)
	if strings.Contains(symbol, ".") {
		return errors.New("cannot use symbol with '.'")
	} else if symbol == "" {
		return errors.New("cannot use empty symbol")
	} else if strings.ContainsAny(symbol, " \n\t") {
		return errors.New("cannot have whitespace in symbol")
	}

	s.mu.Lock()
	s.scope[symbol] = v
	s.mu.Unlock()
	return nil
}

func (s *mapEnv) Resolve(symbol string) (Value, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, found := s.scope[symbol]
	if !found {
		if s.parent == nil {
			return nil, fmt.Errorf("%w: %v", errors.New("unable to resolve symbol"), symbol)
		}
		return s.Parent().Resolve(symbol)
	}

	return v, nil
}

func (s *mapEnv) Parent() Env { return s.parent }

// NewEvalErr returns EvalError with appropriate context added.
func NewEvalErr(v Value, err error) EvalError {
	if ee, ok := err.(EvalError); ok {
		return ee
	} else if ee, ok := err.(*EvalError); ok && ee != nil {
		return *ee
	}

	return EvalError{
		Position: getPosition(v),
		Cause:    err,
		Form:     v,
	}
}

// EvalError represents error during evaluation.
type EvalError struct {
	Position
	Cause error
	Form  Value
}

// Unwrap returns the underlying cause of this error.
func (ee EvalError) Unwrap() error { return ee.Cause }

func (ee EvalError) Error() string {
	return fmt.Sprintf("eval-error in '%s' (at line %d:%d): %v",
		ee.File, ee.Line, ee.Column, ee.Cause,
	)
}

func getPosition(form Value) Position {
	p, hasPosition := form.(interface {
		GetPos() (file string, line, col int)
	})
	if !hasPosition {
		return Position{}
	}

	file, line, col := p.GetPos()
	return Position{
		File:   file,
		Line:   line,
		Column: col,
	}
}
