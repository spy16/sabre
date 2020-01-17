package sabre

import (
	"fmt"
	"sync"
)

// NewScope returns an instance of MapScope with no bindings. If includeCore
// is true, core functions like def, fn, eval etc. will be bound in the new
// scope.
func NewScope(parent Scope, includeCore bool) *MapScope {
	scope := &MapScope{
		parent:   parent,
		mu:       new(sync.RWMutex),
		bindings: map[string]Value{},
	}

	if includeCore {
		_ = scope.Bind("fn", Fn(Lambda))
		_ = scope.Bind("def", Fn(Def))
		_ = scope.Bind("eval", Fn(evalFn))
	}

	return scope
}

// MapScope implements Scope using a Go native hash-map.
type MapScope struct {
	parent   Scope
	mu       *sync.RWMutex
	bindings map[string]Value
}

// BindGo is similar to Bind but handles covnertion of Go value 'v' to
// sabre Val type.
func (scope *MapScope) BindGo(symbol string, v interface{}) error {
	return scope.Bind(symbol, ValueOf(v))
}

// Bind adds the given value to the scope and binds the symbol to it.
func (scope *MapScope) Bind(symbol string, v Value) error {
	scope.mu.Lock()
	defer scope.mu.Unlock()

	scope.bindings[symbol] = v
	return nil
}

// Resolve finds the value bound to the given symbol and returns it if
// found in this scope or parent scope if any.
func (scope *MapScope) Resolve(symbol string) (Value, error) {
	scope.mu.RLock()
	defer scope.mu.RUnlock()

	v, found := scope.bindings[symbol]
	if !found {
		if scope.parent != nil {
			return scope.parent.Resolve(symbol)
		}

		return nil, fmt.Errorf("unable to resolve symbol: %v", symbol)
	}

	return v, nil
}
