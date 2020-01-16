package sabre_test

import (
	"reflect"
	"testing"

	"github.com/spy16/sabre"
)

func TestString_Eval(t *testing.T) {
	executeEvalTests(t, []evalTestCase{
		{
			name:     "Success",
			getScope: nil,
			value:    sabre.String("hello"),
			want:     sabre.String("hello"),
		},
	})
}

func TestKeyword_Eval(t *testing.T) {
	executeEvalTests(t, []evalTestCase{
		{
			name:     "Success",
			getScope: nil,
			value:    sabre.Keyword("hello"),
			want:     sabre.Keyword("hello"),
		},
	})
}

func TestSymbol_Eval(t *testing.T) {
	executeEvalTests(t, []evalTestCase{
		{
			name: "Success",
			getScope: func() sabre.Scope {
				scope := sabre.NewScope(nil)
				scope.Bind("hello", sabre.String("world"))

				return scope
			},
			value: sabre.Symbol("hello"),
			want:  sabre.String("world"),
		},
	})
}

func TestCharacter_Eval(t *testing.T) {
	executeEvalTests(t, []evalTestCase{
		{
			name:     "Success",
			getScope: nil,
			value:    sabre.Character('a'),
			want:     sabre.Character('a'),
		},
	})
}

type evalTestCase struct {
	name     string
	getScope func() sabre.Scope
	value    sabre.Value
	want     sabre.Value
	wantErr  bool
}

func executeEvalTests(t *testing.T, tests []evalTestCase) {
	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var scope sabre.Scope
			if tt.getScope != nil {
				scope = tt.getScope()
			}

			got, err := tt.value.Eval(scope)
			if (err != nil) != tt.wantErr {
				t.Errorf("Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Eval() got = %v, want %v", got, tt.want)
			}
		})
	}
}
