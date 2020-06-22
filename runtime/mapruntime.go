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
		specials: map[string]specialForm{
			"do":    doForm,
			"def":   defForm,
			"cond":  condForm,
			"quote": quoteForm,
		},
	}
}

type mapRuntime struct {
	mu       sync.RWMutex
	scope    map[string]Value
	parent   Runtime
	specials map[string]specialForm
}

type specialForm func(rt Runtime, args ...Value) (specialInvoke, error)

type specialInvoke func() (Value, error)

func (rt *mapRuntime) Eval(form Value) (Value, error) {
	if isNil(form) {
		return Nil{}, nil
	}

	if err := rt.analyze(form); err != nil {
		return nil, err
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

func (rt *mapRuntime) analyze(form Value) error {
	if list, isList := form.(*linkedList); isList {
		return rt.analyzeList(list)
	}
	if seq, ok := ToSeq(form); !ok {
		var err error
		ForEach(seq, func(item Value) bool {
			err = rt.analyze(item)
			return err != nil // break if error
		})
		return err
	}
	return nil
}

func (rt *mapRuntime) analyzeList(list *linkedList) (err error) {
	if list.Count() == 0 {
		return nil
	}

	sym, isSymbol := list.First().(Symbol)
	if !isSymbol {
		return nil
	}

	special, found := rt.specials[sym.Value]
	if !found {
		return nil
	}

	list.specialInvoke, err = special(rt, toSlice(list.Next())...)
	return err
}

func toSlice(seq Seq) []Value {
	var slice []Value
	ForEach(seq, func(item Value) bool {
		slice = append(slice, item)
		return false
	})
	return slice
}
