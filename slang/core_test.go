package slang_test

import (
	"reflect"
	"testing"

	"github.com/spy16/sabre"
	"github.com/spy16/sabre/slang"
)

func TestEval(t *testing.T) {
	t.Parallel()

	table := []struct {
		name     string
		getScope func() sabre.Scope
		arg      sabre.Value
		want     sabre.Value
		wantErr  bool
	}{
		{
			name: "Simple",
			arg:  sabre.Int64(10),
			want: sabre.Int64(10),
		},
		{
			name:     "EvalFailed",
			getScope: func() sabre.Scope { return sabre.NewScope(nil) },
			arg:      sabre.Symbol{Value: "hello"},
			wantErr:  true,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			var scope sabre.Scope
			if tt.getScope != nil {
				scope = tt.getScope()
			}

			got, err := slang.Eval(scope, tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Eval() error = %#v, wantErr %#v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Eval() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestNot(t *testing.T) {
	t.Parallel()

	table := []struct {
		name string
		arg  sabre.Value
		want sabre.Value
	}{
		{
			name: "TruthyValue",
			arg:  sabre.String("hello"),
			want: sabre.Bool(false),
		},
		{
			name: "FalsyValue",
			arg:  sabre.Bool(false),
			want: sabre.Bool(true),
		},
		{
			name: "NoValue",
			arg:  nil,
			want: sabre.Bool(true),
		},
		{
			name: "Nil",
			arg:  sabre.Nil{},
			want: sabre.Bool(true),
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := slang.Not(tt.arg)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %#v, want %#v", got, tt.want)
			}
		})
	}
}
