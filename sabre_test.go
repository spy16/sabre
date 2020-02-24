package sabre_test

import (
	"reflect"
	"testing"

	"github.com/spy16/sabre"
)

const sampleProgram = `
(def v [1 2 3])
(def pi 3.1412)
(def echo (fn* [arg] arg))
(echo pi)

(def int-num 10)
(def float-num 10.1234)
(def list '(nil 1 []))
(def vector ["hello" nil])
(def set #{1 2 3})
(def empty-set #{})

(def complex-calc (let* [sample '(1 2 3 4 [])]
					((. First sample))))

(assert (= int-num 10)
		(= float-num 10.1234)
		(= pi 3.1412)
		(= list '(nil 1 []))
		(= vector ["hello" nil])
		(= empty-set #{})
		(= echo (fn* [arg] arg))
		(= complex-calc 1))

(echo pi)
`

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
			src:  sampleProgram,
			want: sabre.Float64(3.1412),
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			scope := sabre.Scope(sabre.New())
			if tt.getScope != nil {
				scope = tt.getScope()
			}

			scope.Bind("=", sabre.ValueOf(sabre.Compare))
			scope.Bind("assert", &sabre.Fn{Func: assert(t)})

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

func assert(t *testing.T) func(sabre.Scope, []sabre.Value) (sabre.Value, error) {
	return func(scope sabre.Scope, exprs []sabre.Value) (sabre.Value, error) {
		var res sabre.Value
		var err error

		for _, expr := range exprs {
			res, err = expr.Eval(scope)
			if err != nil {
				t.Errorf("%s: %s", expr, err)
			}

			if !isTruthy(res) {
				t.Errorf("assertion failed: %s (result=%v)", expr, res)
			}
		}

		return res, err
	}
}

func isTruthy(v sabre.Value) bool {
	if v == nil || v == (sabre.Nil{}) {
		return false
	}
	if b, ok := v.(sabre.Bool); ok {
		return bool(b)
	}
	return true
}
