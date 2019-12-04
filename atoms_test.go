package sabre_test

import (
	"reflect"
	"testing"

	"github.com/spy16/sabre"
)

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

func TestString_Eval(t *testing.T) {
	executeEvalTests(t, []evalTestCase{
		{
			name:     "Success",
			getScope: nil,
			value:    sabre.String("hello"),
			want:     sabre.String("hello"),
		},
	})
}

func TestKeyword_Eval(t *testing.T) {
	executeEvalTests(t, []evalTestCase{
		{
			name:     "Success",
			getScope: nil,
			value:    sabre.Keyword("hello"),
			want:     sabre.Keyword("hello"),
		},
	})
}

func TestSymbol_Eval(t *testing.T) {
	executeEvalTests(t, []evalTestCase{
		{
			name: "Success",
			getScope: func() sabre.Scope {
				scope := sabre.NewScope(nil)
				scope.Bind("hello", sabre.String("world"))

				return scope
			},
			value: sabre.Symbol("hello"),
			want:  sabre.String("world"),
		},
	})
}

func TestCharacter_Eval(t *testing.T) {
	executeEvalTests(t, []evalTestCase{
		{
			name:     "Success",
			getScope: nil,
			value:    sabre.Character('a'),
			want:     sabre.Character('a'),
		},
	})
}

type evalTestCase struct {
	name     string
	getScope func() sabre.Scope
	value    sabre.Value
	want     sabre.Value
	wantErr  bool
}

func executeEvalTests(t *testing.T, tests []evalTestCase) {
	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var scope sabre.Scope
			if tt.getScope != nil {
				scope = tt.getScope()
			}

			got, err := tt.value.Eval(scope)
			if (err != nil) != tt.wantErr {
				t.Errorf("Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Eval() got = %v, want %v", got, tt.want)
			}
		})
	}
}
