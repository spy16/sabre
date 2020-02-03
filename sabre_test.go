package sabre_test

import (
	"reflect"
	"testing"

	"github.com/spy16/sabre"
)

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
