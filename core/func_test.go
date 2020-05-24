package core

import (
	"reflect"
	"testing"
)

var (
	_ Value      = (*MultiFn)(nil)
	_ Invokable  = (*MultiFn)(nil)
	_ Comparable = (*MultiFn)(nil)

	_ Value      = (*Fn)(nil)
	_ Invokable  = (*Fn)(nil)
	_ Comparable = (*Fn)(nil)
)

func TestMultiFn_Eval(t *testing.T) {
	t.Parallel()

	actual := MultiFn{}

	got, err := New(nil).Eval(actual)
	if err != nil {
		t.Errorf("Eval() unexpected error: %+v", err)
	}

	if !Compare(actual, got) {
		t.Errorf("Eval() want=%+v, got=%+v", actual, got)
	}
}

func TestMultiFn_Invoke(t *testing.T) {
	t.Parallel()

	table := []struct {
		name    string
		getEnv  func() Env
		multiFn MultiFn
		args    []Value
		want    Value
		wantErr bool
	}{
		{
			name: "WrongArity",
			multiFn: MultiFn{
				Name: "arityOne",
				Methods: []Fn{
					{
						Args: []string{"arg1"},
					},
				},
			},
			args:    []Value{},
			wantErr: true,
		},
		{
			name: "VariadicArity",
			multiFn: MultiFn{
				Name: "arityMany",
				Methods: []Fn{
					{
						Args:     []string{"args"},
						Variadic: true,
					},
				},
			},
			args: []Value{},
			want: Nil{},
		},
		{
			name:   "ArgEvalFailure",
			getEnv: func() Env { return New(nil) },
			multiFn: MultiFn{
				Name: "arityOne",
				Methods: []Fn{
					{
						Args: []string{"arg1"},
					},
				},
			},
			args:    []Value{Symbol{Value: "argVal"}},
			wantErr: true,
		},
		{
			name: "Macro",
			getEnv: func() Env {
				env := New(nil)
				env.Bind("argVal", String("hello"))
				return env
			},
			multiFn: MultiFn{
				Name:    "arityOne",
				IsMacro: true,
				Methods: []Fn{
					{
						Args: []string{"arg1"},
						Body: Int64(10),
					},
				},
			},
			args: []Value{Symbol{Value: "argVal"}},
			want: Int64(10),
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			var env Env
			if tt.getEnv != nil {
				env = tt.getEnv()
			}

			got, err := tt.multiFn.Invoke(env, tt.args...)
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
		name    string
		getEnv  func() Env
		fn      Fn
		args    []Value
		want    Value
		wantErr bool
	}{
		{
			name: "GoFuncWrap",
			fn: Fn{
				Func: func(_ Env, _ []Value) (Value, error) {
					return Int64(10), nil
				},
			},
			want: Int64(10),
		},
		{
			name: "NoBody",
			fn: Fn{
				Args: []string{"test"},
			},
			args: []Value{Bool(true)},
			want: Nil{},
		},
		{
			name: "VariadicMatch",
			fn: Fn{
				Args:     []string{"test"},
				Variadic: true,
			},
			args: []Value{},
			want: Nil{},
		},
		{
			name: "VariadicMatch",
			fn: Fn{
				Args:     []string{"test"},
				Variadic: true,
			},
			args: []Value{Int64(10), Bool(true)},
			want: Nil{},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			var env Env
			if tt.getEnv != nil {
				env = tt.getEnv()
			}

			got, err := tt.fn.Invoke(env, tt.args...)
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
