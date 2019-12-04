package sabre_test

import (
	"testing"

	"github.com/spy16/sabre"
)

func TestReader_One_Vector(t *testing.T) {
	executeAllReaderTests(t, []readerTestCase{
		{
			name: "Empty",
			src:  `[]`,
			want: sabre.Vector(nil),
		},
		{
			name: "WithOneEntry",
			src:  `[help]`,
			want: sabre.Vector([]sabre.Value{
				sabre.Symbol("help"),
			}),
		},
		{
			name: "WithMultipleEntry",
			src:  `[+ 0xF 3.1413]`,
			want: sabre.Vector([]sabre.Value{
				sabre.Symbol("+"),
				sabre.Int64(15),
				sabre.Float64(3.1413),
			}),
		},
		{
			name: "WithCommaSeparator",
			src:  `[+,0xF,3.1413]`,
			want: sabre.Vector([]sabre.Value{
				sabre.Symbol("+"),
				sabre.Int64(15),
				sabre.Float64(3.1413),
			}),
		},
		{
			name: "MultiLine",
			src: `[+
                      0xF
                      3.1413
					]`,
			want: sabre.Vector([]sabre.Value{
				sabre.Symbol("+"),
				sabre.Int64(15),
				sabre.Float64(3.1413),
			}),
		},
		{
			name: "MultiLineWithComments",
			src: `[+         ; plus operator adds numerical values
                      0xF    ; hex representation of 15
                      3.1413 ; value of math constant pi
                  ]`,
			want: sabre.Vector([]sabre.Value{
				sabre.Symbol("+"),
				sabre.Int64(15),
				sabre.Float64(3.1413),
			}),
		},
		{
			name:    "UnexpectedEOF",
			src:     "[+ 1 2 ",
			wantErr: true,
		},
	})
}
