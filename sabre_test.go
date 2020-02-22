package sabre_test

import (
	"reflect"
	"testing"

	"github.com/spy16/sabre"
)

func BenchmarkEval(b *testing.B) {
	scope := sabre.NewScope(nil)
	_ = scope.BindGo("inc", func(a int) int {
		return a + 1
	})

	f := &sabre.List{
		Values: sabre.Values{
			sabre.Symbol{Value: "inc"},
			sabre.Int64(10),
		},
	}

	for i := 0; i < b.N; i++ {
		_, _ = sabre.Eval(scope, f)
	}
}

func BenchmarkGoCall(b *testing.B) {
	inc := func(a int) int {
		return a + 1
	}

	for i := 0; i < b.N; i++ {
		_ = inc(10)
	}
}

func TestEval(t *testing.T) {
	t.Parallel()

	table := []struct {
		name     string
		src      string
		getScope func() sabre.Scope
		want     sabre.Value
		wantErr  bool
	}{
		{
			name: "Empty",
			src:  "",
			want: sabre.Nil{},
		},
		{
			name: "SingleForm",
			src:  "123",
			want: sabre.Int64(123),
		},
		{
			name: "MultiForm",
			src:  `123 [] ()`,
			want: &sabre.List{},
		},
		{
			name: "WithFunctionCalls",
			getScope: func() sabre.Scope {
				scope := sabre.NewScope(nil)
				_ = scope.BindGo("ten?", func(i sabre.Int64) bool {
					return i == 10
				})
				return scope
			},
			src:  `(ten? 10)`,
			want: sabre.Bool(true),
		},
		{
			name:    "ReadError",
			src:     `123 [] (`,
			want:    nil,
			wantErr: true,
		},
		{
			name: "Program",
			getScope: func() sabre.Scope {
				scope := sabre.NewScope(nil)
				scope.Bind("def", sabre.Def)
				scope.Bind("fn*", sabre.Lambda)
				return scope
			},
			src:  sampleProgram,
			want: sabre.Float64(3.1412),
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			var scope sabre.Scope
			if tt.getScope != nil {
				scope = tt.getScope()
			}

			got, err := sabre.ReadEvalStr(scope, tt.src)
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

const sampleProgram = `
(def v [1 2 3])

(def pi 3.1412)

(def echo (fn* [arg] arg))

(echo pi)
`
