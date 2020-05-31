package core

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		r        io.Reader
		fileName string
	}{
		{
			name:     "WithStringReader",
			r:        strings.NewReader(":test"),
			fileName: "<string>",
		},
		{
			name:     "WithBytesReader",
			r:        bytes.NewReader([]byte(":test")),
			fileName: "<bytes>",
		},
		{
			name:     "WihFile",
			r:        os.NewFile(0, "test.lisp"),
			fileName: "test.lisp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rd := NewReader(tt.r)
			if rd == nil {
				t.Errorf("New() should return instance of Reader, got nil")
			} else if rd.File != tt.fileName {
				t.Errorf("File = \"%s\", want = \"%s\"", rd.File, tt.name)
			}
		})
	}
}

func TestReader_SetMacro(t *testing.T) {
	t.Run("UnsetDefaultMacro", func(t *testing.T) {
		rd := NewReader(strings.NewReader("~hello"))
		rd.SetMacro('~', nil, false) // remove unquote operator

		var want Value
		want = Symbol{
			Value: "~hello",
			Position: Position{
				File:   "<string>",
				Line:   1,
				Column: 1,
			},
		}

		got, err := rd.One()
		if err != nil {
			t.Errorf("unexpected error: %#v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %+v, want = %+v", got, want)
		}
	})

	t.Run("CustomMacro", func(t *testing.T) {
		rd := NewReader(strings.NewReader("~hello"))
		rd.SetMacro('~', func(rd *Reader, _ rune) (Value, error) {
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

			return String(ru), nil
		}, false) // override unquote operator

		var want Value
		want = String("hello")

		got, err := rd.One()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %+v, want = %+v", got, want)
		}
	})
}

func TestReader_All(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		want    []Value
		wantErr bool
	}{
		{
			name: "ValidLiteralSample",
			src:  `123 "Hello\tWorld" 12.34 -0xF +010 true nil 0b1010 \a :hello`,
			want: []Value{
				Int64(123),
				String("Hello\tWorld"),
				Float64(12.34),
				Int64(-15),
				Int64(8),
				Bool(true),
				Nil{},
				Int64(10),
				Char('a'),
				Keyword("hello"),
			},
		},
		{
			name: "WithComment",
			src:  `:valid-keyword ; comment should return errSkip`,
			want: []Value{Keyword("valid-keyword")},
		},
		{
			name:    "UnterminatedString",
			src:     `:valid-keyword "unterminated string literal`,
			wantErr: true,
		},
		{
			name: "CommentFollowedByForm",
			src:  `; comment should return errSkip` + "\n" + `:valid-keyword`,
			want: []Value{Keyword("valid-keyword")},
		},
		{
			name:    "UnterminatedList",
			src:     `:valid-keyword (add 1 2`,
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
			got, err := NewReader(strings.NewReader(tt.src)).All()
			if (err != nil) != tt.wantErr {
				t.Errorf("All() error = %#v, wantErr %#v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("All() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestReader_One(t *testing.T) {
	executeReaderTests(t, []readerTestCase{
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
			want: &List{
				Items: []Value{
					Symbol{Value: "unquote"},
					&List{
						Items: []Value{
							Symbol{
								Value: "x",
								Position: Position{
									File:   "<string>",
									Line:   1,
									Column: 3,
								},
							},
							Int64(3),
						},
						Position: Position{
							File:   "<string>",
							Line:   1,
							Column: 2,
						},
					},
				},
				Position: Position{
					File:   "<string>",
					Column: 1,
					Line:   1,
				},
			},
		},
	})
}

func TestReader_One_Number(t *testing.T) {
	executeReaderTests(t, []readerTestCase{
		{
			name: "NumberWithLeadingSpaces",
			src:  "    +1234",
			want: Int64(1234),
		},
		{
			name: "PositiveInt",
			src:  "+1245",
			want: Int64(1245),
		},
		{
			name: "NegativeInt",
			src:  "-234",
			want: Int64(-234),
		},
		{
			name: "PositiveFloat",
			src:  "+1.334",
			want: Float64(1.334),
		},
		{
			name: "NegativeFloat",
			src:  "-1.334",
			want: Float64(-1.334),
		},
		{
			name: "PositiveHex",
			src:  "0x124",
			want: Int64(0x124),
		},
		{
			name: "NegativeHex",
			src:  "-0x124",
			want: Int64(-0x124),
		},
		{
			name: "PositiveOctal",
			src:  "0123",
			want: Int64(0123),
		},
		{
			name: "NegativeOctal",
			src:  "-0123",
			want: Int64(-0123),
		},
		{
			name: "PositiveBinary",
			src:  "0b10",
			want: Int64(2),
		},
		{
			name: "NegativeBinary",
			src:  "-0b10",
			want: Int64(-2),
		},
		{
			name: "PositiveBase2Radix",
			src:  "2r10",
			want: Int64(2),
		},
		{
			name: "NegativeBase2Radix",
			src:  "-2r10",
			want: Int64(-2),
		},
		{
			name: "PositiveBase4Radix",
			src:  "4r123",
			want: Int64(27),
		},
		{
			name: "NegativeBase4Radix",
			src:  "-4r123",
			want: Int64(-27),
		},
		{
			name: "ScientificSimple",
			src:  "1e10",
			want: Float64(1e10),
		},
		{
			name: "ScientificNegativeExponent",
			src:  "1e-10",
			want: Float64(1e-10),
		},
		{
			name: "ScientificWithDecimal",
			src:  "1.5e10",
			want: Float64(1.5e+10),
		},
		{
			name:    "FloatStartingWith0",
			src:     "012.3",
			want:    Float64(012.3),
			wantErr: false,
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
			name:    "DecimalPointInBinary",
			src:     "0b1.0101",
			wantErr: true,
		},
		{
			name:    "InvalidDigitForOctal",
			src:     "08",
			wantErr: true,
		},
		{
			name:    "IllegalNumberFormat",
			src:     "9.3.2",
			wantErr: true,
		},
	})
}

func TestReader_One_String(t *testing.T) {
	executeReaderTests(t, []readerTestCase{
		{
			name: "SimpleString",
			src:  `"hello"`,
			want: String("hello"),
		},
		{
			name: "EscapeQuote",
			src:  `"double quote is \""`,
			want: String(`double quote is "`),
		},
		{
			name: "EscapeSlash",
			src:  `"hello\\world"`,
			want: String(`hello\world`),
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

func TestReader_One_Keyword(t *testing.T) {
	executeReaderTests(t, []readerTestCase{
		{
			name: "SimpleASCII",
			src:  `:test`,
			want: Keyword("test"),
		},
		{
			name: "LeadingTrailingSpaces",
			src:  "          :test          ",
			want: Keyword("test"),
		},
		{
			name: "SimpleUnicode",
			src:  `:∂`,
			want: Keyword("∂"),
		},
		{
			name: "WithSpecialChars",
			src:  `:this-is-valid?`,
			want: Keyword("this-is-valid?"),
		},
		{
			name: "FollowedByMacroChar",
			src:  `:this-is-valid'hello`,
			want: Keyword("this-is-valid"),
		},
	})
}

func TestReader_One_Character(t *testing.T) {
	executeReaderTests(t, []readerTestCase{
		{
			name: "ASCIILetter",
			src:  `\a`,
			want: Char('a'),
		},
		{
			name: "ASCIIDigit",
			src:  `\1`,
			want: Char('1'),
		},
		{
			name: "Unicode",
			src:  `\∂`,
			want: Char('∂'),
		},
		{
			name: "Newline",
			src:  `\newline`,
			want: Char('\n'),
		},
		{
			name: "FormFeed",
			src:  `\formfeed`,
			want: Char('\f'),
		},
		{
			name: "Unicode",
			src:  `\u00AE`,
			want: Char('®'),
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

func TestReader_One_Symbol(t *testing.T) {
	executeReaderTests(t, []readerTestCase{
		{
			name: "SimpleASCII",
			src:  `hello`,
			want: Symbol{
				Value: "hello",
				Position: Position{
					File:   "<string>",
					Line:   1,
					Column: 1,
				},
			},
		},
		{
			name: "Unicode",
			src:  `find-∂`,
			want: Symbol{
				Value: "find-∂",
				Position: Position{
					File:   "<string>",
					Line:   1,
					Column: 1,
				},
			},
		},
		{
			name: "SingleChar",
			src:  `+`,
			want: Symbol{
				Value: "+",
				Position: Position{
					File:   "<string>",
					Line:   1,
					Column: 1,
				},
			},
		},
	})
}

func TestReader_One_List(t *testing.T) {
	executeReaderTests(t, []readerTestCase{
		{
			name: "EmptyList",
			src:  `()`,
			want: &List{
				Items: nil,
				Position: Position{
					File:   "<string>",
					Line:   1,
					Column: 1,
				},
			},
		},
		{
			name: "ListWithOneEntry",
			src:  `(help)`,
			want: &List{
				Items: []Value{
					Symbol{
						Value: "help",
						Position: Position{
							File:   "<string>",
							Line:   1,
							Column: 2,
						},
					},
				},
				Position: Position{
					File:   "<string>",
					Line:   1,
					Column: 1,
				},
			},
		},
		{
			name: "ListWithMultipleEntry",
			src:  `(+ 0xF 3.1413)`,
			want: &List{
				Items: []Value{
					Symbol{
						Value: "+",
						Position: Position{
							File:   "<string>",
							Line:   1,
							Column: 2,
						},
					},
					Int64(15),
					Float64(3.1413),
				},
				Position: Position{
					File:   "<string>",
					Line:   1,
					Column: 1,
				},
			},
		},
		{
			name: "ListWithCommaSeparator",
			src:  `(+,0xF,3.1413)`,
			want: &List{
				Items: []Value{
					Symbol{
						Value: "+",
						Position: Position{
							File:   "<string>",
							Line:   1,
							Column: 2,
						},
					},
					Int64(15),
					Float64(3.1413),
				},
				Position: Position{
					File:   "<string>",
					Line:   1,
					Column: 1,
				},
			},
		},
		{
			name: "MultiLine",
			src: `(+
                      0xF
                      3.1413
					)`,
			want: &List{
				Items: []Value{
					Symbol{
						Value: "+",
						Position: Position{
							File:   "<string>",
							Line:   1,
							Column: 2,
						},
					},
					Int64(15),
					Float64(3.1413),
				},
				Position: Position{
					File:   "<string>",
					Line:   1,
					Column: 1,
				},
			},
		},
		{
			name: "MultiLineWithComments",
			src: `(+         ; plus operator adds numerical values
                      0xF    ; hex representation of 15
                      3.1413 ; value of math constant pi
                  )`,
			want: &List{
				Items: []Value{
					Symbol{
						Value: "+",
						Position: Position{
							File:   "<string>",
							Line:   1,
							Column: 2,
						},
					},
					Int64(15),
					Float64(3.1413),
				},
				Position: Position{
					File:   "<string>",
					Line:   1,
					Column: 1,
				},
			},
		},
		{
			name:    "UnexpectedEOF",
			src:     "(+ 1 2 ",
			wantErr: true,
		},
	})
}

type readerTestCase struct {
	name    string
	src     string
	want    Value
	wantErr bool
}

func executeReaderTests(t *testing.T, tests []readerTestCase) {
	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewReader(strings.NewReader(tt.src)).One()
			if (err != nil) != tt.wantErr {
				t.Errorf("One() error = %#v, wantErr %#v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("One() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}
