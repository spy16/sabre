package sabre_test

import (
	"testing"

	"github.com/spy16/sabre"
)

func TestReader_One_Number(t *testing.T) {
	executeAllReaderTests(t, []readerTestCase{
		{
			name: "NumberWithLeadingSpaces",
			src:  "    +1234",
			want: sabre.Int64(1234),
		},
		{
			name: "PositiveInt",
			src:  "+1245",
			want: sabre.Int64(1245),
		},
		{
			name: "NegativeInt",
			src:  "-234",
			want: sabre.Int64(-234),
		},
		{
			name: "PositiveFloat",
			src:  "+1.334",
			want: sabre.Float64(1.334),
		},
		{
			name: "NegativeFloat",
			src:  "-1.334",
			want: sabre.Float64(-1.334),
		},
		{
			name: "PositiveHex",
			src:  "0x124",
			want: sabre.Int64(0x124),
		},
		{
			name: "NegativeHex",
			src:  "-0x124",
			want: sabre.Int64(-0x124),
		},
		{
			name: "PositiveOctal",
			src:  "0123",
			want: sabre.Int64(0123),
		},
		{
			name: "NegativeOctal",
			src:  "-0123",
			want: sabre.Int64(-0123),
		},
		{
			name: "PositiveBinary",
			src:  "0b10",
			want: sabre.Int64(2),
		},
		{
			name: "NegativeBinary",
			src:  "-0b10",
			want: sabre.Int64(-2),
		},
		{
			name: "PositiveBase2Radix",
			src:  "2r10",
			want: sabre.Int64(2),
		},
		{
			name: "NegativeBase2Radix",
			src:  "-2r10",
			want: sabre.Int64(-2),
		},
		{
			name: "PositiveBase4Radix",
			src:  "4r123",
			want: sabre.Int64(27),
		},
		{
			name: "NegativeBase4Radix",
			src:  "-4r123",
			want: sabre.Int64(-27),
		},
		{
			name: "ScientificSimple",
			src:  "1e10",
			want: sabre.Float64(1e10),
		},
		{
			name: "ScientificNegativeExponent",
			src:  "1e-10",
			want: sabre.Float64(1e-10),
		},
		{
			name: "ScientificWithDecimal",
			src:  "1.5e10",
			want: sabre.Float64(1.5e+10),
		},
		{
			name:    "FloatStartingWith0",
			src:     "012.3",
			want:    sabre.Float64(012.3),
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
