package core_test

import (
	"reflect"
	"testing"

	"github.com/spy16/sabre/core"
	"github.com/spy16/sabre/runtime"
)

func TestModule_Eval(t *testing.T) {
	t.Parallel()

	table := []struct {
		title   string
		module  core.Module
		want    runtime.Value
		wantErr bool
	}{
		{
			title:   "NilModule",
			module:  nil,
			want:    runtime.Nil{},
			wantErr: false,
		},
		{
			title:   "EmptyModule",
			module:  core.Module{},
			want:    runtime.Nil{},
			wantErr: false,
		},
		{
			title:   "SingleForm",
			module:  core.Module{runtime.Int64(0)},
			want:    runtime.Int64(0),
			wantErr: false,
		},
		{
			title:   "MultipleForms",
			module:  core.Module{runtime.Int64(0), runtime.Bool(true)},
			want:    runtime.Bool(true),
			wantErr: false,
		},
		{
			title:   "EvalError",
			module:  core.Module{runtime.Symbol{Value: "blah"}},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range table {
		t.Run(tt.title, func(t *testing.T) {
			got, err := tt.module.Eval(runtime.New(nil))
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

func TestModule_String(t *testing.T) {
	m := core.Module{
		runtime.NewSeq(runtime.Int64(0), runtime.Keyword("hello")),
		runtime.Bool(true),
	}

	want := "(do (0 :hello)\n    true)"
	got := m.String()
	if want != got {
		t.Errorf("Module.String() want=`%s`\ngot =`%s`", want, got)
	}
}
