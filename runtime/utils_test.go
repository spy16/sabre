package runtime_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/spy16/sabre/runtime"
)

func TestEquals(t *testing.T) {
	t.Parallel()

	table := []struct {
		v1, v2 runtime.Value
		want   bool
	}{
		{
			v1:   nil,
			v2:   nil,
			want: true,
		},
		{
			v1:   runtime.Nil{},
			v2:   nil,
			want: true,
		},
		{
			v1:   nil,
			v2:   runtime.Nil{},
			want: true,
		},
		{
			v1:   runtime.Bool(true),
			v2:   runtime.Bool(true),
			want: true,
		},
		{
			v1:   runtime.Bool(true),
			v2:   runtime.Bool(false),
			want: false,
		},
		{
			v1:   runtime.Bool(true),
			v2:   runtime.Nil{},
			want: false,
		},
		{
			v1:   runtime.Char('π'),
			v2:   runtime.Nil{},
			want: false,
		},
		{
			v1:   runtime.Char('π'),
			v2:   runtime.Char('π'),
			want: true,
		},
		{
			v1:   runtime.Float64(3.1412),
			v2:   runtime.Float64(3.1412),
			want: true,
		},
		{
			v1:   runtime.Int64(3),
			v2:   runtime.Symbol{Value: "hello"},
			want: false,
		},
		{
			v1:   runtime.Symbol{Value: "hello"},
			v2:   runtime.Int64(3),
			want: false,
		},
		{
			v1:   runtime.Symbol{Value: "hello"},
			v2:   runtime.Symbol{Value: "hello"},
			want: true,
		},
		{
			v1:   runtime.Keyword("specimen"),
			v2:   runtime.String("specimen"),
			want: false,
		},
		{
			v1:   runtime.String("specimen"),
			v2:   runtime.Keyword("specimen"),
			want: false,
		},
		{
			v1: runtime.NewSeq(
				runtime.Float64(10.3),
				runtime.String("sample"),
			),
			v2: runtime.NewSeq(
				runtime.Float64(10.3),
				runtime.String("sample"),
			),
			want: true,
		},
		{
			v1: runtime.NewSeq(
				runtime.Float64(10.3),
				runtime.String("sample"),
			),
			v2: runtime.NewSeq(
				runtime.Float64(10.3),
			),
			want: false,
		},
		{
			v1: runtime.NewSeq(
				runtime.Float64(10.3),
				runtime.String("sample"),
			),
			v2: runtime.NewSeq(
				runtime.Float64(10.3),
				runtime.Keyword("sample"),
			),
			want: false,
		},
		{
			v1: runtime.NewSeq(
				runtime.Float64(10.3),
				runtime.String("sample"),
			),
			v2:   runtime.Nil{},
			want: false,
		},
	}

	for _, tt := range table {
		title := fmt.Sprintf("%s_%s", reflect.TypeOf(tt.v1), reflect.TypeOf(tt.v1))
		t.Run(title, func(t *testing.T) {
			got := runtime.Equals(tt.v1, tt.v2)
			if got != tt.want {
				t.Errorf("Compare('%+v', '%+v') want=%t, got=%t", tt.v1, tt.v2, tt.want, got)
			}
		})
	}
}
