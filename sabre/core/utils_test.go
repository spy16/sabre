package core_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/spy16/sabre/sabre/core"
)

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
