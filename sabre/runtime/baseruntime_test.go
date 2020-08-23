package runtime_test

import (
	"reflect"
	"testing"

	"github.com/spy16/sabre/sabre/runtime"
)

func TestLisp_Eval(t *testing.T) {
	t.Parallel()

	table := []struct {
		name    string
		globals map[string]runtime.Value
		form    runtime.Value
		want    runtime.Value
		wantErr bool
	}{
		{
			name: "nil",
			form: nil,
			want: runtime.Nil{},
		},
		{
			name: "Bool",
			form: runtime.Bool(true),
			want: runtime.Bool(true),
		},
		{
			name: "String",
			form: runtime.String("hello 😎!"),
			want: runtime.String("hello 😎!"),
		},
		{
			name: "Char",
			form: runtime.Char('π'),
			want: runtime.Char('π'),
		},
		{
			name: "Float64",
			form: runtime.Float64(3.1412),
			want: runtime.Float64(3.1412),
		},
		{
			name: "Int64",
			form: runtime.Int64(10),
			want: runtime.Int64(10),
		},
		{
			name: "Symbol_BuiltinTrue",
			form: runtime.Symbol{Value: "true"},
			want: runtime.Bool(true),
		},
		{
			name: "Symbol_BuiltinFalse",
			form: runtime.Symbol{Value: "false"},
			want: runtime.Bool(false),
		},
		{
			name: "Symbol_BuiltinNil",
			form: runtime.Symbol{Value: "nil"},
			want: runtime.Nil{},
		},
		{
			name: "Symbol_CustomBound",
			form: runtime.Symbol{Value: "pi"},
			want: runtime.Float64(3.1412),
			globals: map[string]runtime.Value{
				"pi": runtime.Float64(3.1412),
			},
		},
		{
			name:    "Symbol_Unbound",
			form:    runtime.Symbol{Value: "foo"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "List_Empty",
			form: &runtime.LinkedList{},
			want: &runtime.LinkedList{},
		},
		{
			name: "List_Nil",
			form: (*runtime.LinkedList)(nil),
			want: runtime.NewSeq(),
		},
		{
			name: "List_GoFunc",
			form: runtime.NewSeq(runtime.GoFunc(func(_ runtime.Runtime, _ ...runtime.Value) (runtime.Value, error) {
				return runtime.String("called"), nil
			})),
			want: runtime.String("called"),
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			l := runtime.New(tt.globals)

			got, err := l.Eval(tt.form)
			if (err != nil) != tt.wantErr {
				t.Errorf("Base.Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Base.Eval() got = %v, want %v", got, tt.want)
			}
		})
	}
}
