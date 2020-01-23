package sabre_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/spy16/sabre"
)

func TestList_Eval(t *testing.T) {
	executeEvalTests(t, []evalTestCase{
		{
			name:  "NilList",
			value: sabre.List(nil),
			want:  sabre.List(nil),
		},
		{
			name:  "EmptyList",
			value: sabre.List{},
			want:  sabre.List(nil),
		},
		{
			name:  "Invocation",
			value: sabre.List{sabre.Symbol{Value: "greet"}, sabre.String("Bob")},
			getScope: func() sabre.Scope {
				scope := sabre.NewScope(nil)
				scope.BindGo("greet", func(name sabre.String) string {
					return fmt.Sprintf("Hello %s!", string(name))
				})
				return scope
			},
			want: sabre.String("Hello Bob!"),
		},
		{
			name:    "NonInvokable",
			value:   sabre.List{sabre.Int64(10), sabre.Keyword("hello")},
			wantErr: true,
		},
		{
			name:  "EvalFailure",
			value: sabre.List{sabre.Symbol{Value: "hello"}},
			getScope: func() sabre.Scope {
				return sabre.NewScope(nil)
			},
			wantErr: true,
		},
	})
}

func TestModule_Eval(t *testing.T) {
	executeEvalTests(t, []evalTestCase{
		{
			name:  "NilModule",
			value: sabre.Module(nil),
			want:  sabre.Nil{},
		},
		{
			name:  "EmptyModule",
			value: sabre.Module{},
			want:  sabre.Nil{},
		},
		{
			name:  "SingleForm",
			value: sabre.Module{sabre.Int64(10)},
			want:  sabre.Int64(10),
		},
		{
			name: "MultiForm",
			value: sabre.Module{
				sabre.Int64(10),
				sabre.String("hello"),
			},
			want: sabre.String("hello"),
		},
	})
}

func TestVector_Eval(t *testing.T) {
	executeEvalTests(t, []evalTestCase{
		{
			name:  "NilVector",
			value: sabre.Vector(nil),
			want:  sabre.Vector(nil),
		},
		{
			name:  "EmptyVector",
			value: sabre.Vector{},
			want:  sabre.Vector(nil),
		},
		{
			name: "EvalFailure",
			getScope: func() sabre.Scope {
				return sabre.NewScope(nil)
			},
			value:   sabre.Vector{sabre.Symbol{Value: "hello"}},
			wantErr: true,
		},
	})
}

func TestList_String(t *testing.T) {
	executeStringTestCase(t, []stringTestCase{
		{
			value: sabre.List(nil),
			want:  "()",
		},
		{
			value: sabre.List{},
			want:  "()",
		},
		{
			value: sabre.List{sabre.Keyword("hello")},
			want:  "(:hello)",
		},
		{
			value: sabre.List{sabre.Keyword("hello"), sabre.List{}},
			want:  "(:hello ())",
		},
		{
			value: sabre.List{sabre.Symbol{Value: "quote"}, sabre.Symbol{Value: "hello"}},
			want:  "(quote hello)",
		},
		{
			value: sabre.List{sabre.Symbol{Value: "quote"}, sabre.List{sabre.Symbol{Value: "hello"}}},
			want:  "(quote (hello))",
		},
	})
}

func TestVector_String(t *testing.T) {
	executeStringTestCase(t, []stringTestCase{
		{
			value: sabre.Vector(nil),
			want:  "[]",
		},
		{
			value: sabre.Vector{},
			want:  "[]",
		},
		{
			value: sabre.Vector{sabre.Keyword("hello")},
			want:  "[:hello]",
		},
		{
			value: sabre.Vector{sabre.Keyword("hello"), sabre.List{}},
			want:  "[:hello ()]",
		},
	})
}

func TestModule_String(t *testing.T) {
	executeStringTestCase(t, []stringTestCase{
		{
			value: sabre.Module(nil),
			want:  "",
		},
		{
			value: sabre.Module{sabre.Symbol{Value: "hello"}},
			want:  "hello",
		},
		{
			value: sabre.Module{sabre.Symbol{Value: "hello"}, sabre.Keyword("world")},
			want:  "hello\n:world",
		},
	})
}

func TestVector_Invoke(t *testing.T) {
	t.Parallel()

	vector := sabre.Vector{sabre.Keyword("hello")}

	table := []struct {
		name     string
		getScope func() sabre.Scope
		args     []sabre.Value
		want     sabre.Value
		wantErr  bool
	}{
		{
			name:    "NoArgs",
			args:    []sabre.Value{},
			wantErr: true,
		},
		{
			name:    "InvalidIndex",
			args:    []sabre.Value{sabre.Int64(10)},
			wantErr: true,
		},
		{
			name:    "ValidIndex",
			args:    []sabre.Value{sabre.Int64(0)},
			want:    sabre.Keyword("hello"),
			wantErr: false,
		},
		{
			name:    "NonIntegerArg",
			args:    []sabre.Value{sabre.Keyword("h")},
			wantErr: true,
		},
		{
			name: "EvalFailure",
			getScope: func() sabre.Scope {
				return sabre.NewScope(nil)
			},
			args:    []sabre.Value{sabre.Symbol{Value: "hello"}},
			wantErr: true,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			var scope sabre.Scope
			if tt.getScope != nil {
				scope = tt.getScope()
			}

			got, err := vector.Invoke(scope, tt.args...)
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
