package sabre

import (
	"reflect"
	"testing"

	"github.com/spy16/sabre/runtime"
)

func TestMultiFn_Eval(t *testing.T) {
	m := MultiFn{
		Name: "hello",
		Functions: []runtime.Fn{
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
	m := MultiFn{
		Name: "hello",
		Functions: []runtime.Fn{
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

	want := "(defn hello\n  (fn [a & b] (do 0.123450)))"
	got := m.String()
	if want != got {
		t.Errorf("MultiFn.String() want=`%s`, got=`%s`", want, got)
	}
}

func TestMultiFn_Equals(t *testing.T) {
	m := MultiFn{
		Name: "hello",
		Functions: []runtime.Fn{
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
		multiFn    MultiFn
		args       []runtime.Value
		want       runtime.Value
		wantErr    bool
	}{
		{
			name: "WrongArity",
			multiFn: MultiFn{
				Name: "arityOne",
				Functions: []runtime.Fn{
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
			multiFn: MultiFn{
				Name: "arityMany",
				Functions: []runtime.Fn{
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
			multiFn: MultiFn{
				Name: "arityOne",
				Functions: []runtime.Fn{
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
