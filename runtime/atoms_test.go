package runtime_test

import (
	"reflect"
	"testing"

	"github.com/spy16/sabre/runtime"
)

func TestKeyword_Invoke(t *testing.T) {
	t.Parallel()
	table := []struct {
		title   string
		getEnv  func() runtime.Runtime
		args    []runtime.Value
		want    runtime.Value
		wantErr bool
	}{
		{
			title:   "ArityMismatch",
			args:    nil,
			wantErr: true,
		},
		{
			title:   "ArgEvalError",
			args:    []runtime.Value{runtime.Symbol{Value: "test"}},
			getEnv:  func() runtime.Runtime { return runtime.New(nil) },
			wantErr: true,
		},
		{
			title:  "NotMap",
			args:   []runtime.Value{runtime.Seq(nil)},
			getEnv: func() runtime.Runtime { return runtime.New(nil) },
			want:   runtime.Nil{},
		},
		{
			title: "WithDefault",
			args: []runtime.Value{
				&fakeMap{
					EntryAtFunc: func(key runtime.Value) runtime.Value {
						return nil
					},
				},
				runtime.Float64(10),
			},
			getEnv: func() runtime.Runtime { return runtime.New(nil) },
			want:   runtime.Float64(10),
		},
		{
			title: "WithoutDefault",
			args: []runtime.Value{
				&fakeMap{EntryAtFunc: func(key runtime.Value) runtime.Value {
					if runtime.Equals(key, runtime.Keyword("specimen")) {
						return runtime.Float64(10)
					}
					return nil
				}},
				runtime.Float64(10),
			},
			getEnv: func() runtime.Runtime { return runtime.New(nil) },
			want:   runtime.Float64(10),
		},
	}

	for _, tt := range table {
		t.Run(tt.title, func(t *testing.T) {
			var env runtime.Runtime
			if tt.getEnv != nil {
				env = tt.getEnv()
			}

			got, err := runtime.Keyword("specimen").Invoke(env, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Keyword.Invoke() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keyword.Invoke() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Eval(t *testing.T) {
	t.Parallel()
	runEvalTests(t, []evalTestCase{
		{
			title: "Nil",
			form:  runtime.Nil{},
			want:  runtime.Nil{},
		},
		{
			title: "Bool",
			form:  runtime.Bool(false),
			want:  runtime.Bool(false),
		},
		{
			title: "Float64",
			form:  runtime.Float64(0.123456789),
			want:  runtime.Float64(0.123456789),
		},
		{
			title: "Int64",
			form:  runtime.Int64(10),
			want:  runtime.Int64(10),
		},
		{
			title: "Char",
			form:  runtime.Char('c'),
			want:  runtime.Char('c'),
		},
		{
			title: "Keyword",
			form:  runtime.Keyword("specimen"),
			want:  runtime.Keyword("specimen"),
		},
		{
			title: "String",
			form:  runtime.String("specimen"),
			want:  runtime.String("specimen"),
		},
		{
			title: "Symbol",
			getEnv: func() runtime.Runtime {
				env := runtime.New(nil)
				_ = env.Bind("π", runtime.Float64(3.1412))
				return env
			},
			form: runtime.Symbol{
				Value:    "π",
				Position: runtime.Position{File: "lisp"},
			},
			want: runtime.Float64(3.1412),
		},
	})
}

func Test_String(t *testing.T) {
	t.Parallel()

	table := []struct {
		val  runtime.Value
		want string
	}{
		{
			val:  runtime.Nil{},
			want: "nil",
		},
		{
			val:  runtime.Bool(true),
			want: "true",
		},
		{
			val:  runtime.Bool(false),
			want: "false",
		},
		{
			val:  runtime.Int64(100),
			want: "100",
		},
		{
			val:  runtime.Int64(-100),
			want: "-100",
		},
		{
			val:  runtime.Float64(0.123456),
			want: "0.123456",
		},
		{
			val:  runtime.Float64(-0.123456),
			want: "-0.123456",
		},
		{
			val:  runtime.Float64(0.12345678),
			want: "0.123457",
		},
		{
			val:  runtime.Char('π'),
			want: `\π`,
		},
		{
			val:  runtime.Keyword("specimen"),
			want: ":specimen",
		},
		{
			val:  runtime.Symbol{Value: "specimen"},
			want: "specimen",
		},
		{
			val:  runtime.String("hello 😎"),
			want: `"hello 😎"`,
		},
	}

	for _, tt := range table {
		title := reflect.TypeOf(tt.val).String()
		t.Run(title, func(t *testing.T) {
			got := tt.val.String()
			if tt.want != got {
				t.Errorf("%s.String() want=`%s`, got=`%s`", title, tt.want, got)
			}
		})
	}
}

type fakeMap struct {
	runtime.Map

	EntryAtFunc func(key runtime.Value) runtime.Value
}

func (m *fakeMap) Seq() runtime.Seq {
	return runtime.NewSeq()
}

func (m *fakeMap) Eval(rt runtime.Runtime) (runtime.Value, error) {
	return m, nil
}

func (m *fakeMap) EntryAt(key runtime.Value) runtime.Value {
	return m.EntryAtFunc(key)
}
