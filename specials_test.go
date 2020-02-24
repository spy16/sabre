package sabre_test

import (
	"reflect"
	"testing"

	"github.com/spy16/sabre"
)

func TestSpecials(t *testing.T) {
	scope := sabre.NewScope(nil)
	scope.Bind("def", sabre.Def)
	scope.Bind("let*", sabre.Let)
	scope.Bind("fn*", sabre.Lambda)

	expected := sabre.MultiFn{
		Name:    "hello",
		IsMacro: false,
		Methods: []sabre.Fn{
			{
				Args:     []string{"arg", "rest"},
				Variadic: true,
				Body: sabre.Module{
					sabre.Symbol{Value: "rest"},
				},
			},
		},
	}

	res, err := sabre.ReadEvalStr(scope, src)
	if err != nil {
		t.Errorf("Eval() unexpected error: %v", err)
	}
	if reflect.DeepEqual(res, expected) {
		t.Errorf("Eval() expected=%v, got=%v", expected, res)
	}
}

const src = `
(def temp (let* [pi 3.1412]
			pi))

(def hello (fn* hello
	([arg] arg)
	([arg & rest] rest)))
`
