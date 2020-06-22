package core_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/spy16/sabre/core"
	"github.com/spy16/sabre/runtime"
)

func TestVerifyArgCount(t *testing.T) {
	t.Parallel()

	table := []struct {
		title   string
		argC    int
		arities []int
		wantErr error
	}{
		{
			title:   "ExtraArgs",
			argC:    1,
			arities: []int{},
			wantErr: errors.New("call requires no arguments, got 1"),
		},
		{
			title:   "InsufficientArgs",
			argC:    0,
			arities: []int{1},
			wantErr: errors.New("call requires exactly 1 argument(s), got 0"),
		},
		{
			title:   "ArgCountMismatch",
			argC:    0,
			arities: []int{1, 5},
			wantErr: errors.New("call requires 1 or 5 argument(s), got 0"),
		},
		{
			title:   "ManyArities",
			argC:    4,
			arities: []int{0, 1, 2, 3, 5},
			wantErr: errors.New("wrong number of arguments (4) passed"),
		},
		{
			title:   "Success",
			argC:    2,
			arities: []int{1, 2, 3, 5},
			wantErr: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.title, func(t *testing.T) {
			err := core.VerifyArgCount(tt.arities, tt.argC)
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("VerifyArgCount('%+v', %d) expecting error '%s', got nil",
						tt.arities, tt.argC, tt.wantErr)
				} else if tt.wantErr.Error() != err.Error() {
					t.Errorf("VerifyArgCount('%+v', %d) want=%s, got=%s",
						tt.arities, tt.argC, tt.wantErr, err)
				}
			}
		})
	}
}

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
