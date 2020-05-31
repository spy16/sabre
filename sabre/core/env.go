package core

import (
	"sync"
)

// New returns an empty Env with given parent env. Returned env does not
// support qualified symbol resolution.
func New(parent Env) Env {
	return &mapEnv{
		scope:  map[string]Value{},
		parent: parent,
	}
}

type mapEnv struct {
	mu     sync.RWMutex
	scope  map[string]Value
	parent Env
}

func (env *mapEnv) Eval(form Value) (Value, error) {
	if form == nil {
		return Nil{}, nil
	}

	v, err := form.Eval(env)
	if err != nil {
		e := NewErr(false, getPosition(form), err)
		e.Form = form
		return nil, e
	}

	if v == nil {
		return Nil{}, nil
	}

	return v, nil
}

func (env *mapEnv) Bind(symbol string, v Value) error {
	env.mu.Lock()
	defer env.mu.Unlock()

	env.scope[symbol] = v
	return nil
}

func (env *mapEnv) Resolve(symbol string) (Value, error) {
	env.mu.RLock()
	defer env.mu.RUnlock()

	v, found := env.scope[symbol]
	if !found {
		if env.parent == nil {
			return nil, ErrNotFound
		}

		return env.parent.Resolve(symbol)
	}

	return v, nil
}

func (env *mapEnv) Parent() Env { return env.parent }
