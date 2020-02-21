package sabre

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

var specialForms = map[string]specialForm{}

func init() {
	specialForms = map[string]specialForm{
		"λ":            lambdaForm,
		"fn*":          lambdaForm,
		"if":           ifForm,
		"do":           doForm,
		"def":          defForm,
		"let*":         letForm,
		"throw":        throwErr,
		"quote":        simpleQuote,
		"syntax-quote": syntaxQuote,
	}
}

// lambdaForm defines an anonymous function and returns. Must have the form
// (fn name? [arg*] expr*) or (fn name? ([arg]* expr*)+)
func lambdaForm(scope Scope, args []Value) (specialExpr, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("insufficient args (%d) for 'fn'", len(args))
	}

	def := MultiFn{}
	nextIndex := 0

	name, isName := args[nextIndex].(Symbol)
	if isName {
		def.Name = name.String()
		nextIndex++
	}

	_, isList := args[nextIndex].(*List)
	if isList {
		for _, arg := range args[nextIndex:] {
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
		fn, err := makeFn(scope, args[nextIndex:])
		if err != nil {
			return nil, err
		}
		def.Methods = append(def.Methods, *fn)
	}

	if err := def.validate(); err != nil {
		return nil, err
	}

	return func(_ Scope) (Value, error) {
		return def, nil
	}, nil
}

// ifForm implements if-conditional flow using (if test then else?) form.
func ifForm(scope Scope, args []Value) (specialExpr, error) {
	if err := verifyArgCount([]int{2, 3}, args); err != nil {
		return nil, err
	}

	if err := analyzeSeq(scope, Values(args)); err != nil {
		return nil, err
	}

	return func(scope Scope) (Value, error) {
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
	}, nil
}

// doForm implements the (do <expr>*) special form.
func doForm(scope Scope, args []Value) (specialExpr, error) {
	mod := Module(args)
	if err := analyze(scope, mod); err != nil {
		return nil, err
	}

	return func(scope Scope) (Value, error) {
		return mod.Eval(scope)
	}, nil
}

// defForm implements (def symbol value).
func defForm(scope Scope, args []Value) (specialExpr, error) {
	if err := verifyArgCount([]int{2}, args); err != nil {
		return nil, err
	}

	sym, isSymbol := args[0].(Symbol)
	if !isSymbol {
		return nil, fmt.Errorf("first argument must be symbol, not '%v'",
			reflect.TypeOf(args[0]))
	}

	if err := analyze(scope, args[1]); err != nil {
		return nil, err
	}

	return func(scope Scope) (Value, error) {
		v, err := args[1].Eval(scope)
		if err != nil {
			return nil, err
		}

		if err := rootScope(scope).Bind(sym.String(), v); err != nil {
			return nil, err
		}

		return sym, nil
	}, nil
}

// letForm implements the (let [binding*] expr*) form. expr are evaluated
// with given local bindings.
func letForm(scope Scope, args []Value) (specialExpr, error) {
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

	return func(scope Scope) (Value, error) {
		letScope := NewScope(scope)
		for _, b := range bindings {
			v, err := b.Expr.Eval(letScope)
			if err != nil {
				return nil, err
			}

			_ = letScope.Bind(b.Name, v)
		}

		return Module(args[1:]).Eval(letScope)
	}, nil
}

// throwErr signals an error. Stringified versions of args will be
// concatenated and used as error message.
func throwErr(scope Scope, args []Value) (specialExpr, error) {
	if err := analyzeSeq(scope, Values(args)); err != nil {
		return nil, err
	}

	return func(scope Scope) (Value, error) {
		vals, err := evalValueList(scope, args)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(string(stringFromVals(vals)))
	}, nil
}

// simpleQuote prevents a form from being evaluated.
func simpleQuote(scope Scope, forms []Value) (specialExpr, error) {
	if err := verifyArgCount([]int{1}, forms); err != nil {
		return nil, err
	}

	return func(scope Scope) (Value, error) {
		return forms[0], nil
	}, nil
}

// syntaxQuote recursively applies the quoting to the form.
func syntaxQuote(scope Scope, forms []Value) (specialExpr, error) {
	if err := verifyArgCount([]int{1}, forms); err != nil {
		return nil, err
	}

	return func(scope Scope) (Value, error) {
		quoteScope := NewScope(scope)
		quoteScope.Bind("unquote", Fn{
			Args: []string{"expr"},
			Func: func(scope Scope, args []Value) (Value, error) {
				if err := verifyArgCount([]int{1}, forms); err != nil {
					return nil, err
				}

				return forms[0].Eval(scope)
			},
		})

		return recursiveQuote(quoteScope, forms[0])
	}, nil
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
	if v == nilValue {
		return false
	}

	if b, ok := v.(Bool); ok {
		return bool(b)
	}

	return true
}

func recursiveQuote(scope Scope, f Value) (Value, error) {
	switch v := f.(type) {
	case *List:
		if isUnquote(v.Values) {
			return f.Eval(scope)
		}

		quoted, err := quoteList(scope, v.Values)
		return &List{Values: quoted}, err

	case Set:
		quoted, err := quoteList(scope, v.Values)
		return Set{Values: quoted}, err

	case Vector:
		quoted, err := quoteList(scope, v.Values)
		return Vector{Values: quoted}, err

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

func quoteList(scope Scope, forms []Value) ([]Value, error) {
	var quoted []Value
	for _, form := range forms {
		q, err := recursiveQuote(scope, form)
		if err != nil {
			return nil, err
		}

		quoted = append(quoted, q)
	}

	return quoted, nil
}

func stringFromVals(vals []Value) String {
	argc := len(vals)
	switch argc {
	case 0:
		return String("")

	case 1:
		return String(strings.Trim(vals[0].String(), "\""))

	default:
		var sb strings.Builder
		for _, v := range vals {
			sb.WriteString(strings.Trim(v.String(), "\""))
		}
		return String(sb.String())
	}
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

func analyze(scope Scope, v Value) (err error) {
	switch val := v.(type) {
	case Module:
		err = analyzeSeq(scope, Values(val))

	case *List:
		err = val.parseSpecial(scope)

	case Seq:
		err = analyzeSeq(scope, val)

	}

	if err != nil {
		return EvalError{
			Cause:    err,
			Position: getPosition(v),
			Form:     v,
		}
	}

	return nil
}

func analyzeSeq(scope Scope, seq Seq) error {
	for seq != nil {
		item := seq.First()
		if item == nil {
			break
		}

		if err := analyze(scope, item); err != nil {
			return err
		}

		seq = seq.Next()
	}

	return nil
}

type specialForm func(scope Scope, args []Value) (specialExpr, error)

type specialExpr func(scope Scope) (Value, error)

type binding struct {
	Name string
	Expr Value
}
