package sabre_test

import (
	"testing"

	"github.com/spy16/sabre"
)

func TestList_Eval(t *testing.T) {
	executeEvalTests(t, []evalTestCase{
		{
			name:  "NilList",
			value: sabre.List(nil),
			want:  sabre.List(nil),
		},
		{
			name:  "EmptyList",
			value: sabre.List{},
			want:  sabre.List(nil),
		},
		{
			name:  "QuoteList",
			value: sabre.List{sabre.Symbol("quote"), sabre.Symbol("hello")},
			want:  sabre.Symbol("hello"),
		},
	})
}

func TestModule_Eval(t *testing.T) {
	executeEvalTests(t, []evalTestCase{
		{
			name:  "NilModule",
			value: sabre.Module(nil),
			want:  sabre.Nil{},
		},
		{
			name:  "EmptyModule",
			value: sabre.Module{},
			want:  sabre.Nil{},
		},
		{
			name:  "SingleForm",
			value: sabre.Module{sabre.Int64(10)},
			want:  sabre.Int64(10),
		},
		{
			name: "MultiForm",
			value: sabre.Module{
				sabre.Int64(10),
				sabre.String("hello"),
			},
			want: sabre.String("hello"),
		},
	})
}
