package sabre_test

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/spy16/sabre"
)

func TestNew(t *testing.T) {
	rd := sabre.New(t.Name(), strings.NewReader("test"))
	if rd == nil {
		t.Errorf("New() should return instance of Reader, got nil")
	}
}

func TestReader_SetMacro(t *testing.T) {
	t.Run("UnsetDefaultMacro", func(t *testing.T) {
		rd := sabre.New(t.Name(), strings.NewReader("~hello"))
		rd.SetMacro('~', nil) // remove unquote operator

		var want sabre.Form
		want = sabre.Symbol("~hello")

		got, err := rd.One()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want = %v", got, want)
		}
	})

	t.Run("CustomMacro", func(t *testing.T) {
		rd := sabre.New(t.Name(), strings.NewReader("~hello"))
		rd.SetMacro('~', func(rd *sabre.Reader, _ rune) (sabre.Form, error) {
			var ru []rune
			for {
				r, err := rd.NextRune()
				if err != nil {
					if err == io.EOF {
						break
					}
					return nil, err
				}

				if rd.IsTerminal(r) {
					break
				}
				ru = append(ru, r)
			}

			return sabre.String(ru), nil
		}) // override unquote operator

		var want sabre.Form
		want = sabre.String("hello")

		got, err := rd.One()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want = %v", got, want)
		}
	})
}

func TestReader_All(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		want    sabre.Module
		wantErr bool
	}{
		{
			name: "ValidLiteralSample",
			src:  `"hello\tworld" 123 12.34 -0xF   +010    0b1010     \a :hello 'hello`,
			want: sabre.Module{
				sabre.String("hello\tworld"),
				sabre.Number{Int: 123},
				sabre.Number{IsFloat: true, Float: 12.34},
				sabre.Number{Int: -15},
				sabre.Number{Int: 8},
				sabre.Number{Int: 10},
				sabre.Character('a'),
				sabre.Keyword(":hello"),
				sabre.List{Forms: []sabre.Form{
					sabre.Symbol("quote"),
					sabre.Symbol("hello"),
				}},
			},
		},
		{
			name: "WithComment",
			src:  `:valid-keyword ; comment should return ErrSkip`,
			want: sabre.Module{sabre.Keyword(":valid-keyword")},
		},
		{
			name:    "UnterminatedString",
			src:     `:valid-keyword "unterminated string literal`,
			wantErr: true,
		},
		{
			name: "CommentFollowedByForm",
			src:  `; comment should return ErrSkip` + "\n" + `:valid-keyword`,
			want: sabre.Module{sabre.Keyword(":valid-keyword")},
		},
		{
			name:    "UnterminatedList",
			src:     `:valid-keyword (add 1 2`,
			wantErr: true,
		},
		{
			name:    "UnterminatedVector",
			src:     `:valid-keyword [1 2`,
			wantErr: true,
		},
		{
			name:    "EOFAfterQuote",
			src:     `:valid-keyword '`,
			wantErr: true,
		},
		{
			name:    "CommentAfterQuote",
			src:     `:valid-keyword ';hello world`,
			wantErr: true,
		},
		{
			name:    "UnbalancedParenthesis",
			src:     `())`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sabre.New(tt.name, strings.NewReader(tt.src)).All()
			if (err != nil) != tt.wantErr {
				t.Errorf("All() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("All() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReader_One_Number(t *testing.T) {
	executeAllReaderTests(t, []readerTestCase{
		{
			name: "NumberWithLeadingSpaces",
			src:  "    +1234",
			want: sabre.Number{Int: 1234},
		},
		{
			name: "PositiveInt",
			src:  "+1245",
			want: sabre.Number{Int: 1245},
		},
		{
			name: "NegativeInt",
			src:  "-234",
			want: sabre.Number{Int: -234},
		},
		{
			name: "PositiveFloat",
			src:  "+1.334",
			want: sabre.Number{IsFloat: true, Float: 1.334},
		},
		{
			name: "NegativeFloat",
			src:  "-1.334",
			want: sabre.Number{IsFloat: true, Float: -1.334},
		},
		{
			name: "PositiveHex",
			src:  "0x124",
			want: sabre.Number{Int: 292},
		},
		{
			name: "NegativeHex",
			src:  "-0x124",
			want: sabre.Number{Int: -292},
		},
		{
			name: "PositiveOctal",
			src:  "0123",
			want: sabre.Number{Int: 83},
		},
		{
			name: "NegativeOctal",
			src:  "0123",
			want: sabre.Number{Int: 83},
		},
		{
			name: "PositiveBinary",
			src:  "0b10",
			want: sabre.Number{Int: 2},
		},
		{
			name: "NegativeBinary",
			src:  "-0b10",
			want: sabre.Number{Int: -2},
		},
		{
			name: "PositiveBase2Radix",
			src:  "2r10",
			want: sabre.Number{Int: 2},
		},
		{
			name: "NegativeBase2Radix",
			src:  "-2r10",
			want: sabre.Number{Int: -2},
		},
		{
			name: "PositiveBase4Radix",
			src:  "4r123",
			want: sabre.Number{Int: 27},
		},
		{
			name: "NegativeBase4Radix",
			src:  "-4r123",
			want: sabre.Number{Int: -27},
		},
		{
			name: "ScientificSimple",
			src:  "1e10",
			want: sabre.Number{IsFloat: true, Float: 1e10},
		},
		{
			name: "ScientificNegativeExponent",
			src:  "1e-10",
			want: sabre.Number{IsFloat: true, Float: 1e-10},
		},
		{
			name: "ScientificWithDecimal",
			src:  "1.5e10",
			want: sabre.Number{IsFloat: true, Float: 1.5e10},
		},
		{
			name:    "InvalidValue",
			src:     "1ABe13",
			wantErr: true,
		},
		{
			name:    "InvalidScientificFormat",
			src:     "1e13e10",
			wantErr: true,
		},
		{
			name:    "InvalidExponent",
			src:     "1e1.3",
			wantErr: true,
		},
		{
			name:    "InvalidRadixFormat",
			src:     "1r2r3",
			wantErr: true,
		},
		{
			name:    "RadixBase3WithDigit4",
			src:     "-3r1234",
			wantErr: true,
		},
		{
			name:    "RadixMissingValue",
			src:     "2r",
			wantErr: true,
		},
		{
			name:    "RadixInvalidBase",
			src:     "2ar",
			wantErr: true,
		},
		{
			name:    "RadixWithFloat",
			src:     "2.3r4",
			wantErr: true,
		},
		{
			name:    "FloatStartingWith0",
			src:     "012.3",
			want:    sabre.Number{IsFloat: true, Float: 12.3},
			wantErr: false,
		},
		{
			name:    "DecimalPointInBinary",
			src:     "0b1.0101",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "InvalidDigitForOctal",
			src:     "08",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "IllegalNumberFormat",
			src:     "9.3.2",
			want:    nil,
			wantErr: true,
		},
	})
}

func TestReader_One_Keyword(t *testing.T) {
	executeAllReaderTests(t, []readerTestCase{
		{
			name: "SimpleASCII",
			src:  `:test`,
			want: sabre.Keyword(":test"),
		},
		{
			name: "LeadingTrailingSpaces",
			src:  "          :test          ",
			want: sabre.Keyword(":test"),
		},
		{
			name: "SimpleUnicode",
			src:  `:∂`,
			want: sabre.Keyword(":∂"),
		},
		{
			name: "WithSpecialChars",
			src:  `:this-is-valid?`,
			want: sabre.Keyword(":this-is-valid?"),
		},
		{
			name: "FollowedByMacroChar",
			src:  `:this-is-valid'hello`,
			want: sabre.Keyword(":this-is-valid"),
		},
	})
}

func TestReader_One_Character(t *testing.T) {
	executeAllReaderTests(t, []readerTestCase{
		{
			name: "ASCIILetter",
			src:  `\a`,
			want: sabre.Character('a'),
		},
		{
			name: "ASCIIDigit",
			src:  `\1`,
			want: sabre.Character('1'),
		},
		{
			name: "Unicode",
			src:  `\∂`,
			want: sabre.Character('∂'),
		},
		{
			name: "Newline",
			src:  `\newline`,
			want: sabre.Character('\n'),
		},
		{
			name: "FormFeed",
			src:  `\formfeed`,
			want: sabre.Character('\f'),
		},
		{
			name: "Unicode",
			src:  `\u00AE`,
			want: sabre.Character('®'),
		},
		{
			name:    "InvalidUnicode",
			src:     `\uHELLO`,
			wantErr: true,
		},
		{
			name:    "OutOfRangeUnicode",
			src:     `\u-100`,
			wantErr: true,
		},
		{
			name:    "UnknownSpecial",
			src:     `\hello`,
			wantErr: true,
		},
		{
			name:    "EOF",
			src:     `\`,
			wantErr: true,
		},
	})
}

func TestReader_One_String(t *testing.T) {
	executeAllReaderTests(t, []readerTestCase{
		{
			name: "SimpleString",
			src:  `"hello"`,
			want: sabre.String("hello"),
		},
		{
			name: "EscapeQuote",
			src:  `"double quote is \""`,
			want: sabre.String(`double quote is "`),
		},
		{
			name: "EscapeSlash",
			src:  `"hello\\world"`,
			want: sabre.String(`hello\world`),
		},
		{
			name:    "UnexpectedEOF",
			src:     `"double quote is`,
			wantErr: true,
		},
		{
			name:    "InvalidEscape",
			src:     `"hello \x world"`,
			wantErr: true,
		},
		{
			name:    "EscapeEOF",
			src:     `"hello\`,
			wantErr: true,
		},
	})
}

func TestReader_One_Symbol(t *testing.T) {
	executeAllReaderTests(t, []readerTestCase{
		{
			name: "SimpleASCII",
			src:  `hello`,
			want: sabre.Symbol("hello"),
		},
		{
			name: "Unicode",
			src:  `find-∂`,
			want: sabre.Symbol("find-∂"),
		},
		{
			name: "SingleChar",
			src:  `+`,
			want: sabre.Symbol("+"),
		},
	})
}

func TestReader_One_List(t *testing.T) {
	executeAllReaderTests(t, []readerTestCase{
		{
			name: "EmptyList",
			src:  `()`,
			want: sabre.List{Forms: nil},
		},
		{
			name: "ListWithOneEntry",
			src:  `(help)`,
			want: sabre.List{Forms: []sabre.Form{
				sabre.Symbol("help"),
			}},
		},
		{
			name: "ListWithMultipleEntry",
			src:  `(+ 0xF 3.1413)`,
			want: sabre.List{Forms: []sabre.Form{
				sabre.Symbol("+"),
				sabre.Number{IsFloat: false, Int: 15},
				sabre.Number{IsFloat: true, Float: 3.1413},
			}},
		},
		{
			name: "ListWithCommaSeparator",
			src:  `(+,0xF,3.1413)`,
			want: sabre.List{Forms: []sabre.Form{
				sabre.Symbol("+"),
				sabre.Number{IsFloat: false, Int: 15},
				sabre.Number{IsFloat: true, Float: 3.1413},
			}},
		},
		{
			name: "MultiLine",
			src: `(+ 
                      0xF 
                      3.1413
					)`,
			want: sabre.List{Forms: []sabre.Form{
				sabre.Symbol("+"),
				sabre.Number{IsFloat: false, Int: 15},
				sabre.Number{IsFloat: true, Float: 3.1413},
			}},
		},
		{
			name: "MultiLineWithComments",
			src: `(+         ; plus operator adds numerical values
                      0xF    ; hex representation of 15
                      3.1413 ; value of math constant pi
                  )`,
			want: sabre.List{Forms: []sabre.Form{
				sabre.Symbol("+"),
				sabre.Number{IsFloat: false, Int: 15},
				sabre.Number{IsFloat: true, Float: 3.1413},
			}},
		},
		{
			name:    "UnexpectedEOF",
			src:     "(+ 1 2 ",
			wantErr: true,
		},
	})
}

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
			want: sabre.Vector([]sabre.Form{
				sabre.Symbol("help"),
			}),
		},
		{
			name: "WithMultipleEntry",
			src:  `[+ 0xF 3.1413]`,
			want: sabre.Vector([]sabre.Form{
				sabre.Symbol("+"),
				sabre.Number{IsFloat: false, Int: 15},
				sabre.Number{IsFloat: true, Float: 3.1413},
			}),
		},
		{
			name: "WithCommaSeparator",
			src:  `[+,0xF,3.1413]`,
			want: sabre.Vector([]sabre.Form{
				sabre.Symbol("+"),
				sabre.Number{IsFloat: false, Int: 15},
				sabre.Number{IsFloat: true, Float: 3.1413},
			}),
		},
		{
			name: "MultiLine",
			src: `[+ 
                      0xF 
                      3.1413
					]`,
			want: sabre.Vector([]sabre.Form{
				sabre.Symbol("+"),
				sabre.Number{IsFloat: false, Int: 15},
				sabre.Number{IsFloat: true, Float: 3.1413},
			}),
		},
		{
			name: "MultiLineWithComments",
			src: `[+         ; plus operator adds numerical values
                      0xF    ; hex representation of 15
                      3.1413 ; value of math constant pi
                  ]`,
			want: sabre.Vector([]sabre.Form{
				sabre.Symbol("+"),
				sabre.Number{IsFloat: false, Int: 15},
				sabre.Number{IsFloat: true, Float: 3.1413},
			}),
		},
		{
			name:    "UnexpectedEOF",
			src:     "[+ 1 2 ",
			wantErr: true,
		},
	})
}

type readerTestCase struct {
	name    string
	src     string
	want    sabre.Form
	wantErr bool
}

func executeAllReaderTests(t *testing.T, tests []readerTestCase) {
	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sabre.New(tt.name, strings.NewReader(tt.src)).One()
			if (err != nil) != tt.wantErr {
				t.Errorf("One() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("One() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReader_One(t *testing.T) {
	executeAllReaderTests(t, []readerTestCase{
		{
			name:    "Empty",
			src:     "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "QuotedEOF",
			src:     "';comment is a no-op form\n",
			wantErr: true,
		},
		{
			name:    "ListEOF",
			src:     "( 1",
			wantErr: true,
		},
		{
			name: "UnQuote",
			src:  "~(x 3)",
			want: sabre.List{Forms: []sabre.Form{
				sabre.Symbol("unquote"),
				sabre.List{Forms: []sabre.Form{
					sabre.Symbol("x"),
					sabre.Number{Int: 3},
				}},
			}},
		},
	})
}
