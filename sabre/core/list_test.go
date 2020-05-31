package core_test

import (
	"testing"

	"github.com/spy16/sabre/sabre/core"
)

func TestList_Eval(t *testing.T) {
	t.Parallel()

	runEvalTests(t, []evalTestCase{
		{
			title: "EmptyList",
			form: &core.List{
				Items: []core.Value{},
			},
			want: &core.List{
				Items: []core.Value{},
			},
		},
		{
			title:  "FirstEvalFailure",
			getEnv: func() core.Env { return core.New(nil) },
			form: &core.List{Items: []core.Value{
				core.Symbol{Value: "non-existent"},
			}},
			wantErr: true,
		},
		{
			title:  "NonInvokable",
			getEnv: func() core.Env { return core.New(nil) },
			form: &core.List{Items: []core.Value{
				core.Int64(0),
			}},
			wantErr: true,
		},
		{
			title:  "NonInvokable",
			getEnv: func() core.Env { return core.New(nil) },
			form: &core.List{Items: []core.Value{
				core.GoFunc(func(env core.Env, args ...core.Value) (core.Value, error) {
					return core.String("called"), nil
				}),
			}},
			want:    core.String("called"),
			wantErr: false,
		},
	})
}

func TestList_String(t *testing.T) {
	l := &core.List{
		Items: []core.Value{
			core.Bool(true),
			core.Int64(10),
			core.Float64(3.1412),
		},
	}
	want := `(true 10 3.141200)`
	got := l.String()

	if want != got {
		t.Errorf("List.String() want=`%s`, got=`%s`",
			want, got)
	}
}
