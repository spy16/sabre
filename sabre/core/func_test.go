package core_test

import (
	"reflect"
	"testing"

	"github.com/spy16/sabre/sabre/core"
	"github.com/spy16/sabre/sabre/runtime"
)

func TestMultiFn_Eval(t *testing.T) {
	m := core.MultiFn{
		Name: "hello",
		Functions: []core.Fn{
			{
				Args:     []string{"a", "b"},
				Variadic: true,
				Body:     nil,
			},
		},
	}

	got, err := m.Eval(nil)
	if err != nil {
		t.Errorf("MultiFn.Eval() unexpected error: %+v", err)
	}

	if !reflect.DeepEqual(m, got) {
		t.Errorf("MultiFn.Eval() want=%+v, got=%+v", m, got)
	}
}

func TestMultiFn_String(t *testing.T) {
	m := core.MultiFn{
		Name: "hello",
		Functions: []core.Fn{
			{
				Args:     []string{"a", "b"},
				Variadic: true,
				Body: runtime.NewSeq(
					runtime.Symbol{Value: "do"},
					runtime.Float64(0.12345),
				),
			},
		},
	}

	want := "(defn hello\n  ([a & b] (do 0.123450)))"
	got := m.String()
	if want != got {
		t.Errorf("MultiFn.String() want=`%s`, got=`%s`", want, got)
	}
}

func TestMultiFn_Equals(t *testing.T) {
	m := core.MultiFn{
		Name: "hello",
		Functions: []core.Fn{
			{
				Args:     []string{"a", "b"},
				Variadic: true,
				Body: runtime.NewSeq(
					runtime.Symbol{Value: "do"},
					runtime.Float64(0.12345),
				),
			},
		},
	}

	if !m.Equals(m) {
		t.Errorf("MultiFn.Equals() want=true, got=false")
	}
}

func TestMultiFn_Invoke(t *testing.T) {
	t.Parallel()

	table := []struct {
		name       string
		getRuntime func() runtime.Runtime
		multiFn    core.MultiFn
		args       []runtime.Value
		want       runtime.Value
		wantErr    bool
	}{
		{
			name: "WrongArity",
			multiFn: core.MultiFn{
				Name: "arityOne",
				Functions: []core.Fn{
					{
						Args: []string{"arg1"},
					},
				},
			},
			args:    []runtime.Value{},
			wantErr: true,
		},
		{
			name: "VariadicArity",
			multiFn: core.MultiFn{
				Name: "arityMany",
				Functions: []core.Fn{
					{
						Args:     []string{"args"},
						Variadic: true,
					},
				},
			},
			args: []runtime.Value{},
			want: runtime.Nil{},
		},
		{
			name:       "ArgEvalFailure",
			getRuntime: func() runtime.Runtime { return runtime.New(nil) },
			multiFn: core.MultiFn{
				Name: "arityOne",
				Functions: []core.Fn{
					{
						Args: []string{"arg1"},
					},
				},
			},
			args:    []runtime.Value{runtime.Symbol{Value: "argVal"}},
			wantErr: true,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			var rt runtime.Runtime
			if tt.getRuntime != nil {
				rt = tt.getRuntime()
			}

			got, err := tt.multiFn.Invoke(rt, tt.args...)
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

func TestFn_Eval(t *testing.T) {
	fn := &core.Fn{
		Args:     []string{"a", "b"},
		Variadic: true,
		Body:     runtime.NewSeq(),
	}

	res, err := fn.Eval(nil)
	if err != nil {
		t.Errorf("Fn.Eval() unexpected error: %+v", err)
	}

	if !reflect.DeepEqual(fn, res) {
		t.Errorf("Fn.Eval() want=%+v, got=%+v", fn, res)
	}
}

func TestFn_Invoke(t *testing.T) {
	t.Parallel()
	table := []struct {
		title   string
		fn      *core.Fn
		getRT   func() runtime.Runtime
		args    []runtime.Value
		wantErr bool
		want    runtime.Value
	}{
		{
			title: "NoBody_NoArgs",
			fn:    &core.Fn{},
			want:  runtime.Nil{},
		},
	}

	for _, tt := range table {
		t.Run(tt.title, func(t *testing.T) {
			var rt runtime.Runtime
			if tt.getRT != nil {
				rt = tt.getRT()
			}

			got, err := tt.fn.Invoke(rt)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fn.Invoke() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fn.Invoke() got = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestFn_Equals(t *testing.T) {
	t.Parallel()

	fn := &core.Fn{
		Args:     []string{"a", "b"},
		Variadic: true,
		Body:     runtime.Float64(1.3),
	}

	table := []struct {
		title string
		other runtime.Value
		want  bool
	}{
		{
			title: "SameValue",
			other: fn,
			want:  true,
		},
		{
			title: "DifferentArgs",
			other: &core.Fn{
				Args:     []string{"b"},
				Variadic: true,
				Body:     runtime.Float64(1.3),
			},
			want: false,
		},
		{
			title: "NonVariadic",
			other: &core.Fn{
				Args:     fn.Args,
				Variadic: false,
				Body:     fn.Body,
			},
			want: false,
		},
		{
			title: "DifferentBody",
			other: &core.Fn{
				Args:     fn.Args,
				Variadic: fn.Variadic,
				Body:     runtime.String("something else"),
			},
			want: false,
		},
		{
			title: "NonFnValue",
			other: runtime.Float64(1.),
			want:  false,
		},
	}

	for _, tt := range table {
		t.Run(tt.title, func(t *testing.T) {
			got := fn.Equals(tt.other)
			if tt.want != got {
				t.Errorf("Fn.Equals() want=%+v, got=%+v", tt.want, got)
			}
		})
	}
}

func TestFn_String(t *testing.T) {
	t.Run("Variadic", func(t *testing.T) {
		fn := &core.Fn{
			Args:     []string{"a", "b"},
			Variadic: true,
			Body:     runtime.Float64(1.3),
		}

		want := "(fn [a & b] 1.300000)"
		got := fn.String()

		if want != got {
			t.Errorf("Fn.String() \nwant=`%s`\ngot =`%s`", want, got)
		}
	})

	t.Run("NonVariadic", func(t *testing.T) {
		fn := &core.Fn{
			Args:     []string{"a", "b"},
			Variadic: false,
			Body:     runtime.Float64(1.3),
		}

		want := "(fn [a b] 1.300000)"
		got := fn.String()

		if want != got {
			t.Errorf("Fn.String() \nwant=`%s`\ngot =`%s`", want, got)
		}
	})
}
