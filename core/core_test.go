package core_test

import (
	"reflect"
	"testing"

	"github.com/spy16/sabre"
	"github.com/spy16/sabre/core"
)

func TestCore(t *testing.T) {
	t.Parallel()

	table := []struct {
		name     string
		fn       sabre.Invokable
		args     []sabre.Value
		getScope func() sabre.Scope
		want     sabre.Value
		wantErr  bool
	}{
		{
			name: "Do",
			fn:   core.SpecialFn(core.Do),
			args: []sabre.Value{},
			want: sabre.Nil{},
		},
		{
			name:    "Not_InsufficientArgs",
			fn:      core.Fn(core.Not),
			args:    []sabre.Value{},
			wantErr: true,
		},
		{
			name: "Not_Nil",
			fn:   core.Fn(core.Not),
			args: []sabre.Value{sabre.Nil{}},
			want: sabre.Bool(true),
		},
		{
			name: "Not_Integer",
			fn:   core.Fn(core.Not),
			args: []sabre.Value{sabre.Int64(10)},
			want: sabre.Bool(false),
		},
		{
			name: "Not_False",
			fn:   core.Fn(core.Not),
			args: []sabre.Value{sabre.Bool(false)},
			want: sabre.Bool(true),
		},
		{
			name: "Not_True",
			fn:   core.Fn(core.Not),
			args: []sabre.Value{sabre.Bool(true)},
			want: sabre.Bool(false),
		},
		{
			name: "Def",
			fn:   core.SpecialFn(core.Def),
			args: []sabre.Value{sabre.Symbol("pi"), sabre.Float64(3.1412)},
			getScope: func() sabre.Scope {
				return sabre.NewScope(nil)
			},
			want: sabre.List{sabre.Symbol("quote"), sabre.Symbol("pi")},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			var scope sabre.Scope
			if tt.getScope != nil {
				scope = tt.getScope()
			}

			got, err := tt.fn.Invoke(scope, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLambdaFn(t *testing.T) {
	fn := core.LambdaFn(nil, []sabre.Symbol{"arg1"}, []sabre.Value{sabre.Symbol("arg1")})

	arg1Val := sabre.Int64(10)

	got, err := fn.Invoke(nil, arg1Val)
	if err != nil {
		t.Errorf("Invoke() unexpected error: %v", err)
	}

	if !reflect.DeepEqual(got, arg1Val) {
		t.Errorf("Invoke() want=%#v, got=%#v", arg1Val, got)
	}
}

func TestLambda(t *testing.T) {
	t.Parallel()

	table := []struct {
		name    string
		args    []sabre.Value
		wantErr bool
	}{
		{
			name:    "InsufficientArgs",
			args:    nil,
			wantErr: true,
		},
		{
			name:    "InvalidArgList",
			args:    []sabre.Value{sabre.Int64(0), nil},
			wantErr: true,
		},
		{
			name: "NotSymbolVector",
			args: []sabre.Value{
				sabre.Vector{sabre.Int64(1)},
				sabre.Int64(10),
			},
			wantErr: true,
		},
		{
			name: "Successful",
			args: []sabre.Value{
				sabre.Vector{sabre.Symbol("a"), sabre.Symbol("b")},
				sabre.Int64(10),
			},
			wantErr: false,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got, err := core.Lambda(nil, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Lambda() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("Lambda() expecting non-nil, got nil")
				return
			}
		})
	}
}
