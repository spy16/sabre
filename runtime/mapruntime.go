package runtime

import (
	"sync"
)

// New returns an empty runtime with given parent runtime. Returned runtime does not
// support qualified symbol resolution. parent argument can be nil to make this the
// root runtime.
func New(parent Runtime) Runtime {
	return &mapRuntime{
		scope:  map[string]Value{},
		parent: parent,
	}
}

type mapRuntime struct {
	mu     sync.RWMutex
	scope  map[string]Value
	parent Runtime
}

func (rt *mapRuntime) Eval(form Value) (Value, error) {
	if isNil(form) {
		return Nil{}, nil
	}

	v, err := form.Eval(rt)
	if err != nil {
		e := NewErr(false, getPosition(form), err)
		e.Form = form
		return nil, e
	}

	if isNil(v) {
		return Nil{}, nil
	}
	return v, nil
}

func (rt *mapRuntime) Bind(symbol string, v Value) error {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	rt.scope[symbol] = v
	return nil
}

func (rt *mapRuntime) Resolve(symbol string) (Value, error) {
	rt.mu.RLock()
	defer rt.mu.RUnlock()

	v, found := rt.scope[symbol]
	if !found {
		if rt.parent == nil {
			return nil, ErrNotFound
		}
		return rt.parent.Resolve(symbol)
	}
	return v, nil
}

func (rt *mapRuntime) Parent() Runtime { return rt.parent }

func toSlice(seq Seq) []Value {
	var slice []Value
	ForEach(seq, func(item Value) bool {
		slice = append(slice, item)
		return false
	})
	return slice
}
