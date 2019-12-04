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

func TestReader_One_List(t *testing.T) {
	executeAllReaderTests(t, []readerTestCase{
		{
			name: "EmptyList",
			src:  `()`,
			want: sabre.List(nil),
		},
		{
			name: "ListWithOneEntry",
			src:  `(help)`,
			want: sabre.List{
				sabre.Symbol("help"),
			},
		},
		{
			name: "ListWithMultipleEntry",
			src:  `(+ 0xF 3.1413)`,
			want: sabre.List{
				sabre.Symbol("+"),
				sabre.Int64(15),
				sabre.Float64(3.1413),
			},
		},
		{
			name: "ListWithCommaSeparator",
			src:  `(+,0xF,3.1413)`,
			want: sabre.List{
				sabre.Symbol("+"),
				sabre.Int64(15),
				sabre.Float64(3.1413),
			},
		},
		{
			name: "MultiLine",
			src: `(+
                      0xF
                      3.1413
					)`,
			want: sabre.List{
				sabre.Symbol("+"),
				sabre.Int64(15),
				sabre.Float64(3.1413),
			},
		},
		{
			name: "MultiLineWithComments",
			src: `(+         ; plus operator adds numerical values
                      0xF    ; hex representation of 15
                      3.1413 ; value of math constant pi
                  )`,
			want: sabre.List{
				sabre.Symbol("+"),
				sabre.Int64(15),
				sabre.Float64(3.1413),
			},
		},
		{
			name:    "UnexpectedEOF",
			src:     "(+ 1 2 ",
			wantErr: true,
		},
	})
}
