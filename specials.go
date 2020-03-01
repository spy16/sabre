package sabre

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
)

var (
	// Def implements (def symbol value) form for defining bindings.
	Def = SpecialForm{
		Name:  "def",
		Parse: parseDef,
	}

	// Lambda defines an anonymous function and returns. Must have the form
	// (fn* name? [arg*] expr*) or (fn* name? ([arg]* expr*)+)
	Lambda = SpecialForm{
		Name:  "fn*",
		Parse: fnParser(false),
	}

	// Macro defines an anonymous function and returns. Must have the form
	// (macro* name? [arg*] expr*) or (fn* name? ([arg]* expr*)+)
	Macro = SpecialForm{
		Name:  "macro*",
		Parse: fnParser(true),
	}

	// Let implements the (let [binding*] expr*) form. expr are evaluated
	// with given local bindings.
	Let = SpecialForm{
		Name:  "let",
		Parse: parseLet,
	}

	// Do special form evaluates args one by one and returns the result of
	// the last expr.
	Do = SpecialForm{
		Name:  "do",
		Parse: parseDo,
	}

	// If implements if-conditional flow using (if test then else?) form.
	If = SpecialForm{
		Name:  "if",
		Parse: parseIf,
	}

	// SimpleQuote prevents a form from being evaluated.
	SimpleQuote = SpecialForm{
		Name:  "quote",
		Parse: parseSimpleQuote,
	}

	// SyntaxQuote recursively applies the quoting to the form.
	SyntaxQuote = SpecialForm{
		Name:  "syntax-quote",
		Parse: parseSyntaxQuote,
	}
)

func fnParser(isMacro bool) func(scope Scope, forms []Value) (*Fn, error) {
	return func(scope Scope, forms []Value) (*Fn, error) {
		if len(forms) < 1 {
			return nil, fmt.Errorf("insufficient args (%d) for 'fn'", len(forms))
		}

		nextIndex := 0
		def := MultiFn{
			IsMacro: isMacro,
		}

		name, isName := forms[nextIndex].(Symbol)
		if isName {
			def.Name = name.String()
			nextIndex++
		}

		return &Fn{
			Func: func(_ Scope, args []Value) (Value, error) {
				_, isList := forms[nextIndex].(*List)
				if isList {
					for _, arg := range forms[nextIndex:] {
						spec, isList := arg.(*List)
						if !isList {
							return nil, fmt.Errorf("expected arg to be list, not %s",
								reflect.TypeOf(arg))
						}

						fn, err := makeFn(scope, spec.Values)
						if err != nil {
							return nil, err
						}

						def.Methods = append(def.Methods, *fn)
					}
				} else {
					fn, err := makeFn(scope, forms[nextIndex:])
					if err != nil {
						return nil, err
					}
					def.Methods = append(def.Methods, *fn)
				}
				return def, def.validate()
			},
		}, nil
	}
}

func parseLet(scope Scope, args []Value) (*Fn, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("call requires at-least bindings argument")
	}

	vec, isVector := args[0].(Vector)
	if !isVector {
		return nil, fmt.Errorf(
			"first argument to let must be bindings vector, not %v",
			reflect.TypeOf(args[0]),
		)
	}

	if len(vec.Values)%2 != 0 {
		return nil, fmt.Errorf("bindings must contain event forms")
	}

	var bindings []binding
	for i := 0; i < len(vec.Values); i += 2 {
		sym, isSymbol := vec.Values[i].(Symbol)
		if !isSymbol {
			return nil, fmt.Errorf(
				"item at %d must be symbol, not %s",
				i, vec.Values[i],
			)
		}

		bindings = append(bindings, binding{
			Name: sym.Value,
			Expr: vec.Values[i+1],
		})
	}

	return &Fn{
		Func: func(scope Scope, _ []Value) (Value, error) {
			letScope := NewScope(scope)
			for _, b := range bindings {
				v, err := b.Expr.Eval(letScope)
				if err != nil {
					return nil, err
				}
				_ = letScope.Bind(b.Name, v)
			}
			return Module(args[1:]).Eval(letScope)
		},
	}, nil
}

func parseDo(scope Scope, args []Value) (*Fn, error) {
	return &Fn{
		Func: func(scope Scope, args []Value) (Value, error) {
			if len(args) == 0 {
				return Nil{}, nil
			}

			results, err := evalValueList(scope, args)
			if err != nil {
				return nil, err
			}
			return results[len(results)-1], err
		},
	}, nil
}

func parseDef(scope Scope, forms []Value) (*Fn, error) {
	if err := verifyArgCount([]int{2}, forms); err != nil {
		return nil, err
	}

	if err := analyze(scope, forms[1]); err != nil {
		return nil, err
	}

	return &Fn{
		Func: func(scope Scope, args []Value) (Value, error) {
			sym, isSymbol := args[0].(Symbol)
			if !isSymbol {
				return nil, fmt.Errorf("first argument must be symbol, not '%v'",
					reflect.TypeOf(args[0]))
			}

			v, err := args[1].Eval(scope)
			if err != nil {
				return nil, err
			}

			if err := rootScope(scope).Bind(sym.String(), v); err != nil {
				return nil, err
			}

			return sym, nil
		},
	}, nil
}

func parseIf(scope Scope, args []Value) (*Fn, error) {
	if err := verifyArgCount([]int{2, 3}, args); err != nil {
		return nil, err
	}

	if err := analyzeSeq(scope, Values(args)); err != nil {
		return nil, err
	}

	return &Fn{
		Func: func(scope Scope, args []Value) (Value, error) {
			test, err := args[0].Eval(scope)
			if err != nil {
				return nil, err
			}

			if !isTruthy(test) {
				// handle 'else' flow.
				if len(args) == 2 {
					return Nil{}, nil
				}

				return args[2].Eval(scope)
			}

			// handle 'if true' flow.
			return args[1].Eval(scope)
		},
	}, nil
}

func parseSimpleQuote(scope Scope, forms []Value) (*Fn, error) {
	return &Fn{
		Func: func(scope Scope, _ []Value) (Value, error) {
			return forms[0], verifyArgCount([]int{1}, forms)
		},
	}, nil
}

func parseSyntaxQuote(scope Scope, forms []Value) (*Fn, error) {
	if err := verifyArgCount([]int{1}, forms); err != nil {
		return nil, err
	}

	if err := analyzeSeq(scope, Values(forms)); err != nil {
		return nil, err
	}

	return &Fn{
		Func: func(scope Scope, _ []Value) (Value, error) {
			return recursiveQuote(scope, forms[0])
		},
	}, nil
}

// SpecialForm is a Value type for representing special forms that will be
// subjected to an intermediate Parsing stage before evaluation.
type SpecialForm struct {
	Name  string
	Parse func(scope Scope, args []Value) (*Fn, error)
}

// Eval always returns error since it is not allowed to directly evaluate
// a special form.
func (sf SpecialForm) Eval(_ Scope) (Value, error) {
	return nil, errors.New("can't take value of special form")
}

func (sf SpecialForm) String() string {
	return fmt.Sprintf("SpecialForm{name=%s}", sf.Name)
}

func analyze(scope Scope, form Value) error {
	switch f := form.(type) {
	case Module:
		for _, expr := range f {
			if err := analyze(scope, expr); err != nil {
				return err
			}
		}

	case *List:
		return f.parse(scope)

	case String:
		return nil

	case Seq:
		return analyzeSeq(scope, f)
	}

	return nil
}

func analyzeSeq(scope Scope, seq Seq) error {
	for seq != nil {
		f := seq.First()
		if f == nil {
			break
		}

		if err := analyze(scope, f); err != nil {
			return err
		}
		seq = seq.Next()
	}

	return nil
}

func recursiveQuote(scope Scope, f Value) (Value, error) {
	switch v := f.(type) {
	case *List:
		if isUnquote(v.Values) {
			if err := verifyArgCount([]int{1}, v.Values[1:]); err != nil {
				return nil, err
			}

			return v.Values[1].Eval(scope)
		}

		quoted, err := quoteSeq(scope, v.Values)
		return &List{Values: quoted}, err

	case Set:
		quoted, err := quoteSeq(scope, v.Values)
		return Set{Values: quoted}, err

	case Vector:
		quoted, err := quoteSeq(scope, v.Values)
		return Vector{Values: quoted}, err

	case String:
		return f, nil

	case Seq:
		return quoteSeq(scope, v)

	default:
		return f, nil
	}
}

func isUnquote(list []Value) bool {
	if len(list) == 0 {
		return false
	}

	sym, isSymbol := list[0].(Symbol)
	if !isSymbol {
		return false
	}

	return sym.Value == "unquote"
}

func quoteSeq(scope Scope, seq Seq) (Values, error) {
	var quoted []Value
	for seq != nil {
		f := seq.First()
		if f == nil {
			break
		}

		q, err := recursiveQuote(scope, f)
		if err != nil {
			return nil, err
		}

		quoted = append(quoted, q)
		seq = seq.Next()
	}
	return quoted, nil
}

func verifyArgCount(arities []int, args []Value) error {
	actual := len(args)
	sort.Ints(arities)

	if len(arities) == 0 && actual != 0 {
		return fmt.Errorf("call requires no arguments, got %d", actual)
	}

	L := len(arities)
	switch {
	case L == 1 && actual != arities[0]:
		return fmt.Errorf("call requires exactly %d argument(s), got %d", arities[0], actual)

	case L == 2:
		c1, c2 := arities[0], arities[1]
		if actual != c1 && actual != c2 {
			return fmt.Errorf("call requires %d or %d argument(s), got %d", c1, c2, actual)
		}

	case L > 2:
		return fmt.Errorf("wrong number of arguments (%d) passed", actual)
	}

	return nil
}

func rootScope(scope Scope) Scope {
	if scope == nil {
		return nil
	}
	p := scope
	for temp := scope; temp != nil; temp = temp.Parent() {
		p = temp
	}
	return p
}

func isTruthy(v Value) bool {
	if v == (Nil{}) {
		return false
	}
	if b, ok := v.(Bool); ok {
		return bool(b)
	}
	return true
}

func makeFn(scope Scope, spec []Value) (*Fn, error) {
	if len(spec) < 1 {
		return nil, fmt.Errorf("insufficient args (%d) for 'fn'", len(spec))
	}

	body := Module(spec[1:])
	if err := analyze(scope, body); err != nil {
		return nil, err
	}

	fn := &Fn{Body: body}
	if err := fn.parseArgSpec(spec[0]); err != nil {
		return nil, err
	}

	return fn, nil
}

type binding struct {
	Name string
	Expr Value
}

func accessMember(target reflect.Value, member string) (reflect.Value, error) {
	if member[0] >= 'a' && member[0] <= 'z' {
		return reflect.Value{}, fmt.Errorf("cannot access private member")
	}

	if _, found := target.Type().MethodByName(member); found {
		return target.MethodByName(member), nil
	}

	if target.Kind() == reflect.Ptr {
		target = target.Elem()
	}

	if _, found := target.Type().FieldByName(member); found {
		return target.FieldByName(member), nil
	}

	return reflect.Value{}, fmt.Errorf("value of type '%s' has no member named '%s'",
		target.Type(), member)
}
