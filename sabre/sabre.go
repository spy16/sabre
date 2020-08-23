package sabre

import "errors"

const globalFrame = "<global>"

// ErrNotFound is returned by Sabre when a symbol resolution fails.
var ErrNotFound = errors.New("not found")

// New returns a new instance of Sabre with given options. If no runtime is
// provided, a default runtime will be used.
func New(opts ...Option) *Sabre {
	s := &Sabre{}
	s.push(stackFrame{name: globalFrame})
	for _, opt := range withDefaults(opts) {
		opt(s)
	}
	return s
}

// Sabre represents an instance of sabre interpreter and acts an environment
// for evaluating forms/expressions.
type Sabre struct {
	rt       Runtime
	stack    []stackFrame
	maxDepth int
	specials map[string]ParseSpecial
}

// Eval evaluates the given form and returns the resultant value or error.
func (s *Sabre) Eval(form Value) (Value, error) {
	if seq, isSeq := form.(Seq); isSeq {
		temp, err := s.macroExpand(seq)
		if err != nil {
			return nil, err
		} else if temp != nil {
			form = temp
		}
	}

	expr, err := s.analyze(form)
	if err != nil {
		return nil, err
	} else if expr == nil {
		return Nil{}, nil
	}

	return expr.Eval(s)
}

func (s *Sabre) macroExpand(seq Seq) (Value, error) {
	// TODO: Implement macro expansion.
	return nil, nil // return nil value since no expansion
}

func (s *Sabre) analyze(form Value) (Expr, error) {
	if isNil(form) {
		return &ConstExpr{Value: Nil{}}, nil
	}

	switch v := form.(type) {
	case Symbol:
		val, err := s.resolve(v.Value)
		if err != nil {
			return nil, err
		}
		return &ConstExpr{Value: val}, nil

	case Expr:
		return v, nil

	case Seq:
		if v.Count() == 0 {
			return &ConstExpr{Value: v}, nil
		}
		return s.analyzeSeq(v)
	}

	if analyzer, ok := s.rt.(Analyzer); ok {
		return analyzer.Analyze(s, form)
	}

	return &ConstExpr{Value: form}, nil
}

func (s *Sabre) analyzeSeq(seq Seq) (Expr, error) {
	if expansion, err := s.macroExpand(seq); err != nil {
		return nil, err
	} else if expansion != nil {
		// macro expansion did happen. throw away the sequence and continue
		// with the expanded form.
		return s.analyze(expansion)
	}

	// handle special form analysis.
	if sym, ok := seq.First().(Symbol); ok {
		parse, found := s.specials[sym.Value]
		if found {
			return parse(s, seq.Next())
		}
	}

	return s.parseInvoke(seq)
}

func (s *Sabre) parseInvoke(seq Seq) (*InvokeExpr, error) {
	val, err := s.analyze(seq.First())
	if err != nil {
		return nil, err
	}
	seq = seq.Next()

	var args []Expr
	for ; seq != nil; seq = seq.Next() {
		arg, err := s.analyze(seq.First())
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}

	return &InvokeExpr{
		Target: val,
		Args:   args,
	}, nil
}

func (s *Sabre) push(frame stackFrame) {
	if frame.vars == nil {
		frame.vars = map[string]Value{}
	}
	s.stack = append(s.stack, frame)
}

func (s *Sabre) pop() *stackFrame {
	if len(s.stack) == 0 {
		panic("runtime stack must never be empty")
	}

	f := s.stack[len(s.stack)-1]
	s.stack = s.stack[0 : len(s.stack)-1]
	return &f
}

func (s *Sabre) resolve(name string) (Value, error) {
	if len(s.stack) == 0 {
		panic("runtime stack must never be empty")
	}

	for i := len(s.stack) - 1; i >= 0; i-- {
		if v, found := s.stack[i].vars[name]; found {
			return v, nil
		}
	}

	return nil, ErrNotFound
}

type stackFrame struct {
	name string
	args Seq
	vars map[string]Value

	// positional information
	file      string
	line, col string
}
