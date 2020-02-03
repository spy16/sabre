package sabre_test

import (
	"reflect"
	"testing"

	"github.com/spy16/sabre"
)

func TestMultiFn_Eval(t *testing.T) {
	executeEvalTests(t, []evalTestCase{
		{
			name:  "Valid",
			value: sabre.MultiFn{},
			want:  sabre.MultiFn{},
		},
	})
}

func TestMultiFn_String(t *testing.T) {
	executeStringTestCase(t, []stringTestCase{
		{
			value: sabre.MultiFn{
				Name: "hello",
			},
			want: "MultiFn{name=hello}",
		},
	})
}

func TestMultiFn_Invoke(t *testing.T) {
	t.Parallel()

	table := []struct {
		name     string
		getScope func() sabre.Scope
		multiFn  sabre.MultiFn
		args     []sabre.Value
		want     sabre.Value
		wantErr  bool
	}{
		{
			name: "WrongArity",
			multiFn: sabre.MultiFn{
				Name: "arityOne",
				Methods: []sabre.Fn{
					sabre.Fn{
						Args: []string{"arg1"},
					},
				},
			},
			args:    []sabre.Value{},
			wantErr: true,
		},
		{
			name: "VariadicArity",
			multiFn: sabre.MultiFn{
				Name: "arityMany",
				Methods: []sabre.Fn{
					sabre.Fn{
						Args:     []string{"args"},
						Variadic: true,
					},
				},
			},
			args: []sabre.Value{},
			want: sabre.Nil{},
		},
		{
			name:     "ArgEvalFailure",
			getScope: func() sabre.Scope { return sabre.NewScope(nil) },
			multiFn: sabre.MultiFn{
				Name: "arityOne",
				Methods: []sabre.Fn{
					sabre.Fn{
						Args: []string{"arg1"},
					},
				},
			},
			args:    []sabre.Value{sabre.Symbol{Value: "argVal"}},
			wantErr: true,
		},
		{
			name: "Macro",
			getScope: func() sabre.Scope {
				scope := sabre.NewScope(nil)
				scope.Bind("argVal", sabre.String("hello"))
				return scope
			},
			multiFn: sabre.MultiFn{
				Name:    "arityOne",
				IsMacro: true,
				Methods: []sabre.Fn{
					sabre.Fn{
						Args: []string{"arg1"},
						Body: sabre.Int64(10),
					},
				},
			},
			args: []sabre.Value{sabre.Symbol{Value: "argVal"}},
			want: sabre.Int64(10),
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			var scope sabre.Scope
			if tt.getScope != nil {
				scope = tt.getScope()
			}

			got, err := tt.multiFn.Invoke(scope, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Invoke() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Invoke() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFn_Invoke(t *testing.T) {
	t.Parallel()

	table := []struct {
		name     string
		getScope func() sabre.Scope
		fn       sabre.Fn
		args     []sabre.Value
		want     sabre.Value
		wantErr  bool
	}{
		{
			name: "GoFuncWrap",
			fn: sabre.Fn{
				Func: sabre.GoFunc(func(scope sabre.Scope, args []sabre.Value) (sabre.Value, error) {
					return sabre.Int64(10), nil
				}),
			},
			want: sabre.Int64(10),
		},
		{
			name: "NoBody",
			fn: sabre.Fn{
				Args: []string{"test"},
			},
			args: []sabre.Value{sabre.Bool(true)},
			want: sabre.Nil{},
		},
		{
			name: "VariadicMatch",
			fn: sabre.Fn{
				Args:     []string{"test"},
				Variadic: true,
			},
			args: []sabre.Value{},
			want: sabre.Nil{},
		},
		{
			name: "VariadicMatch",
			fn: sabre.Fn{
				Args:     []string{"test"},
				Variadic: true,
			},
			args: []sabre.Value{sabre.Int64(10), sabre.Bool(true)},
			want: sabre.Nil{},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			var scope sabre.Scope
			if tt.getScope != nil {
				scope = tt.getScope()
			}

			got, err := tt.fn.Invoke(scope, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Invoke() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Invoke() got = %v, want %v", got, tt.want)
			}
		})
	}
}
