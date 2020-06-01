package core_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/spy16/sabre/sabre/core"
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
				}
				if tt.wantErr.Error() != err.Error() {
					t.Errorf("VerifyArgCount('%+v', %d) want=%s, got=%s",
						tt.arities, tt.argC, tt.wantErr, err)
				}
			}
		})
	}
}

func TestCompare(t *testing.T) {
	t.Parallel()

	table := []struct {
		v1, v2 core.Value
		want   bool
	}{
		{
			v1:   nil,
			v2:   nil,
			want: true,
		},
		{
			v1:   core.Nil{},
			v2:   nil,
			want: true,
		},
		{
			v1:   nil,
			v2:   core.Nil{},
			want: true,
		},
		{
			v1:   core.Float64(3.1412),
			v2:   core.Float64(3.1412),
			want: true,
		},
		{
			v1:   core.Int64(3),
			v2:   core.Symbol{Value: "hello"},
			want: false,
		},
		{
			v1:   core.Symbol{Value: "hello"},
			v2:   core.Int64(3),
			want: false,
		},
		{
			v1:   core.Symbol{Value: "hello"},
			v2:   core.Symbol{Value: "hello"},
			want: true,
		},
		{
			v1:   core.Keyword("specimen"),
			v2:   core.String("specimen"),
			want: false,
		},
		{
			v1:   core.String("specimen"),
			v2:   core.Keyword("specimen"),
			want: false,
		},
		{
			v1: &core.List{Items: []core.Value{
				core.Float64(10.3),
				core.String("sample"),
			}},
			v2: &core.List{Items: []core.Value{
				core.Float64(10.3),
				core.String("sample"),
			}},
			want: true,
		},
		{
			v1: &core.List{Items: []core.Value{
				core.Float64(10.3),
				core.String("sample"),
			}},
			v2: &core.List{Items: []core.Value{
				core.Float64(10.3),
			}},
			want: false,
		},
		{
			v1: &core.List{Items: []core.Value{
				core.Float64(10.3),
				core.String("sample"),
			}},
			v2: &core.List{Items: []core.Value{
				core.Float64(10.3),
				core.Keyword("sample"),
			}},
			want: false,
		},
		{
			v1: &core.List{Items: []core.Value{
				core.Float64(10.3),
				core.String("sample"),
			}},
			v2:   core.Nil{},
			want: false,
		},
	}

	for _, tt := range table {
		title := fmt.Sprintf("%s_%s", reflect.TypeOf(tt.v1), reflect.TypeOf(tt.v1))
		t.Run(title, func(t *testing.T) {
			got := core.Compare(tt.v1, tt.v2)
			if got != tt.want {
				t.Errorf("Compare('%+v', '%+v') want=%t, got=%t", tt.v1, tt.v2, tt.want, got)
			}
		})
	}
}

func Test_mapEnv(t *testing.T) {
	parent := core.New(nil)
	_ = parent.Bind("π", core.Float64(3.1412))

	env := core.New(parent)
	_ = env.Bind("message", core.String("Hello World!"))

	t.Run("EvalNil", func(t *testing.T) {
		v, err := env.Eval(nil)
		if err != nil {
			t.Errorf("mapEnv.Resolve(\"message\"): unexpected error: %v", err)
		}
		want := core.Nil{}
		if !core.Compare(v, want) {
			t.Errorf("mapEnv.Resolve(\"message\") want=%+v, got=%+v", want, v)
		}
	})

	t.Run("Resolve", func(t *testing.T) {
		v, err := env.Resolve("message")
		if err != nil {
			t.Errorf("mapEnv.Resolve(\"message\"): unexpected error: %v", err)
		}
		want := core.String("Hello World!")
		if !core.Compare(v, want) {
			t.Errorf("mapEnv.Resolve(\"message\") want=%+v, got=%+v", want, v)
		}
	})

	t.Run("ResolveFromParent", func(t *testing.T) {
		v, err := env.Resolve("π")
		if err != nil {
			t.Errorf("mapEnv.Resolve(\"π\"): unexpected error: %v", err)
		}
		want := core.Float64(3.1412)
		if !core.Compare(v, want) {
			t.Errorf("mapEnv.Resolve(\"π\") want=%+v, got=%+v", want, v)
		}
	})
}
