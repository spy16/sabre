package runtime

import (
	"errors"
	"fmt"
	"reflect"
)

const globalFrame = "<global>"

// New returns a base runtime instance with given globals setup.
func New(factory Factory, globals map[string]Value) *Base {
	rt := &Base{factory: factory}
	rt.push(stackFrame{name: globalFrame, vars: globals})

	vars := rt.stack[0].vars
	vars["true"] = Bool(true)
	vars["false"] = Bool(false)
	vars["nil"] = Nil{}
	return rt
}

// Factory is responsible for constructing collection values.
type Factory interface {
	NewMap() Map
	NewVec() Vector
	NewSet() Set
}

// Base implements a Sabre runtime. It is NOT safe for concurrent use without
// external synchronization. Zero value is NOT safe for use.
type Base struct {
	// MaxDepth represents the maximum allowed stack depth. If not set, stack
	// size is unbound.
	MaxDepth int

	stack   []stackFrame
	factory Factory
}

// Eval evaluates the given form and returns the resultant value. Evaluating
// an atom returns the atom itself. Evaluating symbol results in value bound
// for the symbol. Evaluating a list leads to an invocation. Any other value
// type is returned as is.
func (rt *Base) Eval(form Value) (Value, error) {
	form, err := rt.analyze(form)
	if err != nil {
		return nil, err
	}

	if isNil(form) {
		return Nil{}, nil
	}

	switch v := form.(type) {
	case Int64, Float64, Bool, String, Char, Keyword, GoFunc:
		return v, nil

	case Symbol:
		return rt.resolve(v.Value)

	case Seq:
		if v.Count() == 0 {
			return v, nil
		}
		return rt.doInvoke(v)

	case Vector:
		return rt.evalVector(v)

	case Map:
		return rt.evalMap(v)

	case Set:
		return rt.evalSet(v)

	default:
		return v, ErrNoEval
	}
}

func (rt *Base) analyze(form Value) (Value, error) {
	return form, nil
}

func (rt *Base) evalVector(vec Vector) (Vector, error) {
	if rt.factory == nil {
		return nil, errors.New("support for vectors is disabled")
	}

	items, err := EvalAll(rt, toSlice(vec.Seq()))
	if err != nil {
		return nil, err
	}
	return rt.factory.NewVec().Conj(items...), nil
}

func (rt *Base) evalSet(set Set) (Set, error) {
	if rt.factory == nil {
		return nil, errors.New("support for sets is disabled")
	}

	items, err := EvalAll(rt, toSlice(set.Seq()))
	if err != nil {
		return nil, err
	}
	return rt.factory.NewSet().Conj(items...)
}

func (rt *Base) evalMap(m Map) (Map, error) {
	if rt.factory == nil {
		return nil, errors.New("support for maps is disabled")
	}

	res := rt.factory.NewMap()

	for keys := m.Keys(); keys != nil; keys = keys.Next() {
		k := keys.First()

		key, err := rt.Eval(k)
		if err != nil {
			return nil, err
		}

		val, err := rt.Eval(m.EntryAt(k))
		if err != nil {
			return nil, err
		}

		res, err = res.Assoc(key, val)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (rt *Base) bind(local bool, name string, value Value) {
	if !local {
		// try to find the frame where this name is already bound
		// and re-bind it.
		for i := len(rt.stack) - 1; i >= 0; i-- {
			vars := rt.stack[i].vars
			if _, found := vars[name]; found {
				vars[name] = value
				return
			}
		}
	}

	rt.stack[len(rt.stack)-1].vars[name] = value
}

func (rt *Base) resolve(name string) (Value, error) {
	if len(rt.stack) == 0 {
		panic("runtime stack must never be empty")
	}

	for i := len(rt.stack) - 1; i >= 0; i-- {
		if v, found := rt.stack[i].vars[name]; found {
			return v, nil
		}
	}

	return nil, ErrNotFound
}

func (rt *Base) doInvoke(seq Seq) (Value, error) {
	first, args := seq.First(), seq.Next()

	fn, ok := first.(Invokable)
	if !ok {
		return nil, fmt.Errorf("value of type '%s' is not invokable",
			reflect.TypeOf(first))
	}

	if rt.MaxDepth > 0 && len(rt.stack) >= rt.MaxDepth {
		return nil, fmt.Errorf("maximum stack depth reached")
	}
	rt.push(stackFrame{name: fn.String(), args: args})
	defer rt.pop()

	return fn.Invoke(rt, toSlice(args)...)
}

func (rt *Base) push(frame stackFrame) {
	if frame.vars == nil {
		frame.vars = map[string]Value{}
	}
	rt.stack = append(rt.stack, frame)
}

func (rt *Base) pop() *stackFrame {
	if len(rt.stack) == 0 {
		panic("runtime stack must never be empty")
	}

	f := rt.stack[len(rt.stack)-1]
	rt.stack = rt.stack[0 : len(rt.stack)-1]
	return &f
}

type stackFrame struct {
	name string
	args Seq
	vars map[string]Value
}
