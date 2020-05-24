package core

import (
	"reflect"
	"testing"
)

func Test_mapEnv_Eval(t *testing.T) {
	t.Parallel()

	table := []struct {
		title   string
		form    Value
		want    Value
		wantErr bool
	}{
		{
			title:   "nilForm",
			form:    nil,
			want:    Nil{},
			wantErr: false,
		},
		{
			title:   "Number",
			form:    Float64(10),
			want:    Float64(10),
			wantErr: false,
		},
		{
			title: "Expr",
			form: Vector{Values: Values{
				Float64(10),
				String("hello"),
			}},
			want: Vector{Values: Values{
				Float64(10),
				String("hello"),
			}},
			wantErr: false,
		},
	}

	for _, tt := range table {
		t.Run(tt.title, func(t *testing.T) {
			env := New(nil)
			got, err := env.Eval(tt.form)
			if (err != nil) != tt.wantErr {
				t.Errorf("Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Eval() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_mapEnv_Bind(t *testing.T) {
	t.Parallel()

	table := []struct {
		title   string
		sym     string
		val     Value
		wantErr bool
	}{
		{
			title:   "ValidBind",
			sym:     "foo",
			val:     Float64(10),
			wantErr: false,
		},
		{
			title:   "EmptySymbol",
			sym:     "",
			val:     Float64(10),
			wantErr: true,
		},
		{
			title:   "DottedSymbol",
			sym:     "foo.bar",
			val:     Float64(10),
			wantErr: true,
		},
		{
			title:   "WhitespacedSymbol",
			sym:     "foo bar",
			val:     Float64(10),
			wantErr: true,
		},
	}

	for _, tt := range table {
		t.Run(tt.title, func(t *testing.T) {
			env := New(nil)

			err := env.Bind(tt.sym, tt.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("Bind() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_mapEnv_Resolve(t *testing.T) {
	t.Parallel()

	table := []struct {
		title   string
		setup   func() Env
		symbol  string
		want    Value
		wantErr bool
	}{
		{
			title: "Existing",
			setup: func() Env {
				env := New(nil)
				_ = env.Bind("pi", Float64(3.1412))
				return env
			},
			symbol:  "pi",
			want:    Float64(3.1412),
			wantErr: false,
		},
		{
			title: "NonExistent",
			setup: func() Env {
				env := New(nil)
				return env
			},
			symbol:  "pi",
			want:    nil,
			wantErr: true,
		},
		{
			title: "ExistingInParent",
			setup: func() Env {
				parent := New(nil)
				parent.Bind("pi", Float64(3.1412))
				return New(parent)
			},
			symbol:  "pi",
			want:    Float64(3.1412),
			wantErr: false,
		},
	}

	for _, tt := range table {
		t.Run(tt.title, func(t *testing.T) {
			env := tt.setup()

			got, err := env.Resolve(tt.symbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("Resolve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resolve() got = %+v, want %+v", got, tt.want)
			}
		})
	}
}
