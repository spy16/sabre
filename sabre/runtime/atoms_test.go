package runtime_test

import (
	"reflect"
	"testing"

	"github.com/spy16/sabre/sabre/runtime"
)

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
				return runtime.New(map[string]runtime.Value{
					"π": runtime.Float64(3.1412),
				})
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
