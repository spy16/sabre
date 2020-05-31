package core_test

import (
	"reflect"
	"testing"

	"github.com/spy16/sabre/sabre/collection"
	"github.com/spy16/sabre/sabre/core"
)

func TestKeyword_Invoke(t *testing.T) {
	t.Parallel()
	table := []struct {
		title   string
		getEnv  func() core.Env
		args    []core.Value
		want    core.Value
		wantErr bool
	}{
		{
			title:   "ArityMismatch",
			args:    nil,
			wantErr: true,
		},
		{
			title:   "ArgEvalError",
			args:    []core.Value{core.Symbol{Value: "test"}},
			getEnv:  func() core.Env { return core.New(nil) },
			wantErr: true,
		},
		{
			title:  "NotMap",
			args:   []core.Value{&core.List{}},
			getEnv: func() core.Env { return core.New(nil) },
			want:   core.Nil{},
		},
		{
			title:  "WithDefault",
			args:   []core.Value{&collection.HashMap{}, core.Float64(10)},
			getEnv: func() core.Env { return core.New(nil) },
			want:   core.Float64(10),
		},
		{
			title: "WithoutDefault",
			args:  []core.Value{core.Symbol{Value: "m"}, core.Float64(10)},
			getEnv: func() core.Env {
				env := core.New(nil)
				m := core.Map(&collection.HashMap{})
				m, _ = m.Assoc(core.Keyword("specimen"), core.Float64(10))
				env.Bind("m", m)
				return env
			},
			want: core.Float64(10),
		},
	}

	for _, tt := range table {
		t.Run(tt.title, func(t *testing.T) {
			var env core.Env
			if tt.getEnv != nil {
				env = tt.getEnv()
			}

			got, err := core.Keyword("specimen").Invoke(env, tt.args...)
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
			form:  core.Nil{},
			want:  core.Nil{},
		},
		{
			title: "Bool",
			form:  core.Bool(false),
			want:  core.Bool(false),
		},
		{
			title: "Float64",
			form:  core.Float64(0.123456789),
			want:  core.Float64(0.123456789),
		},
		{
			title: "Int64",
			form:  core.Int64(10),
			want:  core.Int64(10),
		},
		{
			title: "Char",
			form:  core.Char('c'),
			want:  core.Char('c'),
		},
		{
			title: "Keyword",
			form:  core.Keyword("specimen"),
			want:  core.Keyword("specimen"),
		},
		{
			title: "String",
			form:  core.String("specimen"),
			want:  core.String("specimen"),
		},
		{
			title: "Symbol",
			getEnv: func() core.Env {
				env := core.New(nil)
				_ = env.Bind("π", core.Float64(3.1412))
				return env
			},
			form: core.Symbol{
				Value:    "π",
				Position: core.Position{File: "lisp"},
			},
			want: core.Float64(3.1412),
		},
	})
}

func Test_String(t *testing.T) {
	t.Parallel()

	table := []struct {
		val  core.Value
		want string
	}{
		{
			val:  core.Nil{},
			want: "nil",
		},
		{
			val:  core.Bool(true),
			want: "true",
		},
		{
			val:  core.Bool(false),
			want: "false",
		},
		{
			val:  core.Int64(100),
			want: "100",
		},
		{
			val:  core.Int64(-100),
			want: "-100",
		},
		{
			val:  core.Float64(0.123456),
			want: "0.123456",
		},
		{
			val:  core.Float64(-0.123456),
			want: "-0.123456",
		},
		{
			val:  core.Float64(0.12345678),
			want: "0.123457",
		},
		{
			val:  core.Char('π'),
			want: `\π`,
		},
		{
			val:  core.Char('😎'),
			want: `\😎`,
		},
		{
			val:  core.Keyword("specimen"),
			want: ":specimen",
		},
		{
			val:  core.Symbol{Value: "specimen"},
			want: "specimen",
		},
		{
			val:  core.String("hello 😎"),
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

type evalTestCase struct {
	title   string
	getEnv  func() core.Env
	form    core.Value
	want    core.Value
	wantErr bool
}

func runEvalTests(t *testing.T, cases []evalTestCase) {
	for _, tt := range cases {
		t.Run(tt.title, func(t *testing.T) {
			var env core.Env
			if tt.getEnv != nil {
				env = tt.getEnv()
			}

			got, err := tt.form.Eval(env)
			if (err != nil) != tt.wantErr {
				t.Errorf("%s.Eval() error = %v, wantErr %v",
					reflect.TypeOf(tt.form), err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s.Eval() got = %v, want %v",
					reflect.TypeOf(tt.form), got, tt.want)
			}
		})
	}
}
