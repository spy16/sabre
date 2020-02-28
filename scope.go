package sabre

import (
	"errors"
	"fmt"
	"sync"
)

// ErrResolving is returned when a scope implementation fails to resolve
// a binding for given symbol.
var ErrResolving = errors.New("unable to resolve symbol")

// New initializes a new scope with all the core bindings.
func New() *MapScope {
	scope := &MapScope{
		parent:   nil,
		mu:       new(sync.RWMutex),
		bindings: map[string]Value{},
	}

	scope.Bind("macroexpand", ValueOf(func(scope Scope, v Value) (Value, error) {
		f, _, err := MacroExpand(scope, v)
		return f, err
	}))

	scope.Bind("quote", SimpleQuote)
	scope.Bind("syntax-quote", SyntaxQuote)

	scope.Bind("fn*", Lambda)
	scope.Bind("macro*", Macro)
	scope.Bind("let*", Let)
	scope.Bind("if", If)
	scope.Bind("do", Do)
	scope.Bind("def", Def)

	return scope
}

// NewScope returns an instance of MapScope with no bindings. If you need
// builtin special forms, pass result of New() as argument.
func NewScope(parent Scope) *MapScope {
	return &MapScope{
		parent:   parent,
		mu:       new(sync.RWMutex),
		bindings: map[string]Value{},
	}
}

// MapScope implements Scope using a Go native hash-map.
type MapScope struct {
	parent   Scope
	mu       *sync.RWMutex
	bindings map[string]Value
}

// Parent returns the parent scope of this scope.
func (scope *MapScope) Parent() Scope { return scope.parent }

// Bind adds the given value to the scope and binds the symbol to it.
func (scope *MapScope) Bind(symbol string, v Value) error {
	scope.mu.Lock()
	defer scope.mu.Unlock()

	scope.bindings[symbol] = v
	return nil
}

// Resolve finds the value bound to the given symbol and returns it if
// found in this scope or parent scope if any. Returns error otherwise.
func (scope *MapScope) Resolve(symbol string) (Value, error) {
	scope.mu.RLock()
	defer scope.mu.RUnlock()

	v, found := scope.bindings[symbol]
	if !found {
		if scope.parent != nil {
			return scope.parent.Resolve(symbol)
		}

		return nil, fmt.Errorf("%w: %v", ErrResolving, symbol)
	}

	return v, nil
}

// BindGo is similar to Bind but handles conversion of Go value 'v' to
// sabre Value type. See `ValueOf()`
func (scope *MapScope) BindGo(symbol string, v interface{}) error {
	return scope.Bind(symbol, ValueOf(v))
}
