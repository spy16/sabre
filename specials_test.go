package sabre_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/spy16/sabre"
)

const src = `
(def temp (let* [pi 3.1412]
			pi))

(def hello (fn* hello
	([arg] arg)
	([arg & rest] rest)))
`

func TestSpecials(t *testing.T) {
	scope := sabre.New()

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

func TestDot(t *testing.T) {
	t.Parallel()

	table := []struct {
		name    string
		src     string
		want    sabre.Value
		wantErr bool
	}{
		{
			name: "StringFieldAccess",
			src:  "(. Name foo)",
			want: sabre.String("Bob"),
		},
		{
			name: "BoolFieldAccess",
			src:  "(. Enabled foo)",
			want: sabre.Bool(false),
		},
		{
			name: "MethodAccess",
			src:  `((. Bar foo) "Baz")`,
			want: sabre.String("Bar(\"Baz\")"),
		},
		{
			name: "MethodAccessPtr",
			src:  `((. BarPtr foo) "Bob")`,
			want: sabre.String("BarPtr(\"Bob\")"),
		},
		{
			name:    "InvalidNumberOfArgs",
			src:     `(. BarPtr)`,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "EvalFailed",
			src:     `(. BarPtr blah)`,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "NonSymbolArgument",
			src:     `(. "BarPtr" foo)`,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "NonExistentMember",
			src:     `(. Baz foo)`,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "PrivateMember",
			src:     `(. privateMember foo)`,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "InvalidGoSymbol",
			src:     `(. baz-$-foo# foo)`,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			scope := sabre.New()
			scope.BindGo("foo", &Foo{
				Name: "Bob",
			})

			form, err := sabre.NewReader(strings.NewReader(tt.src)).All()
			if err != nil {
				t.Fatalf("failed to read source='%s': %+v", tt.src, err)
			}

			got, err := sabre.Eval(scope, form)
			if (err != nil) != tt.wantErr {
				t.Errorf("Eval() unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("Eval() want=%#v, got=%#v", tt.want, got)
			}
		})
	}
}

// Foo is a dummy type for member access tests.
type Foo struct {
	Name    string
	Enabled bool
}

func (foo *Foo) BarPtr(arg string) string {
	return fmt.Sprintf("BarPtr(\"%s\")", arg)
}

func (foo Foo) Bar(arg string) string {
	return fmt.Sprintf("Bar(\"%s\")", arg)
}
