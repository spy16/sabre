package sabre_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/spy16/sabre"
)

func TestBool_Eval(t *testing.T) {
	executeEvalTests(t, []evalTestCase{
		{
			name:     "Success",
			getScope: nil,
			value:    sabre.Bool(true),
			want:     sabre.Bool(true),
		},
	})
}

func TestNil_Eval(t *testing.T) {
	executeEvalTests(t, []evalTestCase{
		{
			name:     "Success",
			getScope: nil,
			value:    sabre.Nil{},
			want:     sabre.Nil{},
		},
	})
}

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

func TestNil_String(t *testing.T) {
	executeStringTestCase(t, []stringTestCase{
		{
			value: sabre.Nil{},
			want:  "nil",
		},
	})
}

func TestInt64_String(t *testing.T) {
	executeStringTestCase(t, []stringTestCase{
		{
			value: sabre.Int64(10),
			want:  "10",
		},
		{
			value: sabre.Int64(-10),
			want:  "-10",
		},
	})
}

func TestFloat64_String(t *testing.T) {
	executeStringTestCase(t, []stringTestCase{
		{
			value: sabre.Float64(10.3),
			want:  "10.300000",
		},
		{
			value: sabre.Float64(-10.3),
			want:  "-10.300000",
		},
	})
}

func TestBool_String(t *testing.T) {
	executeStringTestCase(t, []stringTestCase{
		{
			value: sabre.Bool(true),
			want:  "true",
		},
		{
			value: sabre.Bool(false),
			want:  "false",
		},
	})
}

func TestKeyword_String(t *testing.T) {
	executeStringTestCase(t, []stringTestCase{
		{
			value: sabre.Keyword("hello"),
			want:  ":hello",
		},
	})
}

func TestSymbol_String(t *testing.T) {
	executeStringTestCase(t, []stringTestCase{
		{
			value: sabre.Symbol("hello"),
			want:  "hello",
		},
	})
}

func TestCharacter_String(t *testing.T) {
	executeStringTestCase(t, []stringTestCase{
		{
			value: sabre.Character('a'),
			want:  "\\a",
		},
	})
}

func TestString_String(t *testing.T) {
	executeStringTestCase(t, []stringTestCase{
		{
			value: sabre.String("hello world"),
			want:  `"hello world"`,
		},
		{
			value: sabre.String("hello\tworld"),
			want: `"hello	world"`,
		},
	})
}

type stringTestCase struct {
	value sabre.Value
	want  string
}

type evalTestCase struct {
	name     string
	getScope func() sabre.Scope
	value    sabre.Value
	want     sabre.Value
	wantErr  bool
}

func executeStringTestCase(t *testing.T, tests []stringTestCase) {
	t.Parallel()

	for _, tt := range tests {
		t.Run(reflect.TypeOf(tt.value).Name(), func(t *testing.T) {
			got := strings.TrimSpace(tt.value.String())
			if got != tt.want {
				t.Errorf("String() got = %v, want %v", got, tt.want)
			}
		})
	}
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
