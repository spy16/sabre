package slang

import (
	"fmt"
	"strings"
	"sync"

	"github.com/spy16/sabre"
)

const (
	nsSeparator = '/'
	defaultNS   = "user"
)

// New returns a new instance of Slang interpreter.
func New() *Slang {
	sl := &Slang{
		mu:       &sync.RWMutex{},
		bindings: map[nsSymbol]sabre.Value{},
	}

	_ = sl.SwitchNS(sabre.Symbol{Value: defaultNS})
	_ = sl.BindGo("ns", sl.SwitchNS)

	return sl
}

// Slang represents an instance of slang interpreter.
type Slang struct {
	mu        *sync.RWMutex
	currentNS string
	bindings  map[nsSymbol]sabre.Value
}

// Eval evalautes the given value in Slang context.
func (slang *Slang) Eval(v sabre.Value) (sabre.Value, error) {
	return sabre.Eval(slang, v)
}

// ReadEval reads the source and evalautes it in Slang context.
func (slang *Slang) ReadEval(src string) (sabre.Value, error) {
	return sabre.ReadEvalStr(slang, src)
}

// Bind binds the given name to the given Value into the slang interpreter
// context.
func (slang *Slang) Bind(symbol string, v sabre.Value) error {
	slang.mu.Lock()
	defer slang.mu.Unlock()

	nsSym, err := slang.splitSymbol(symbol)
	if err != nil {
		return err
	}

	if nsSym.NS != slang.currentNS {
		return fmt.Errorf("cannot to bind outside current namespace")
	}

	slang.bindings[*nsSym] = v
	return nil
}

// Resolve finds the value bound to the given symbol and returns it if
// found in the Slang context and returns it.
func (slang *Slang) Resolve(symbol string) (sabre.Value, error) {
	slang.mu.RLock()
	defer slang.mu.RUnlock()

	if symbol == "ns" {
		symbol = "user/ns"
	}

	nsSym, err := slang.splitSymbol(symbol)
	if err != nil {
		return nil, err
	}

	v, found := slang.bindings[*nsSym]
	if !found {
		return nil, fmt.Errorf("unable to resolve symbol: %v", symbol)
	}

	return v, nil
}

// BindGo is similar to Bind but handles convertion of Go value 'v' to
// sabre Value type.
func (slang *Slang) BindGo(symbol string, v interface{}) error {
	return slang.Bind(symbol, sabre.ValueOf(v))
}

// SwitchNS changes the current namespace to the string value of given symbol.
func (slang *Slang) SwitchNS(sym sabre.Symbol) error {
	slang.mu.Lock()
	slang.currentNS = sym.String()
	slang.mu.Unlock()

	return slang.Bind("*ns*", sym)
}

// Parent always returns nil to represent this is the root scope.
func (slang *Slang) Parent() sabre.Scope {
	return nil
}

// CurrentNS returns the current active namespace.
func (slang *Slang) CurrentNS() string {
	slang.mu.RLock()
	defer slang.mu.RUnlock()

	return slang.currentNS
}

func (slang *Slang) splitSymbol(symbol string) (*nsSymbol, error) {
	parts := strings.Split(symbol, string(nsSeparator))
	if len(parts) < 2 {
		return &nsSymbol{
			NS:   slang.currentNS,
			Name: symbol,
		}, nil
	} else if len(parts) > 2 {
		return nil, fmt.Errorf("invalid qualified symbol: '%s'", symbol)
	}

	return &nsSymbol{
		NS:   parts[0],
		Name: parts[1],
	}, nil
}

type nsSymbol struct {
	NS   string
	Name string
}
