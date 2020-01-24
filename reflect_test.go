package sabre

import (
	"reflect"
	"testing"
)

var simpleFn = func() {}
var simpleFnRV = reflect.ValueOf(simpleFn)

var anyVal = struct{ name string }{}
var anyValRV = reflect.ValueOf(anyVal)

func TestValueOf(t *testing.T) {
	t.Parallel()

	table := []struct {
		name string
		v    interface{}
		want Value
	}{
		{
			name: "int64",
			v:    int64(10),
			want: Int64(10),
		},
		{
			name: "float",
			v:    float32(10.),
			want: Float64(10.),
		},
		{
			name: "uint8",
			v:    uint8('a'),
			want: Character('a'),
		},
		{
			name: "bool",
			v:    true,
			want: Bool(true),
		},
		{
			name: "Value",
			v:    Int64(10),
			want: Int64(10),
		},
		{
			name: "Nil",
			v:    nil,
			want: Nil{},
		},
		{
			name: "Any",
			v:    anyVal,
			want: anyValue{rv: anyValRV},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := ValueOf(tt.v)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValueOf() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_strictFn_Invoke(t *testing.T) {
	t.Parallel()

	table := []struct {
		name    string
		v       interface{}
		args    []Value
		want    Value
		wantErr bool
	}{
		{
			name: "SimpleNoArgNoReturn",
			v:    func() {},
			want: nil,
		},
		{
			name: "SimpleNoArg",
			v:    func() int { return 10 },
			want: Int64(10),
		},
		{
			name: "SimpleNoReturn",
			v:    func(arg Int64) {},
			args: []Value{Int64(10)},
			want: nil,
		},
		{
			name: "SimpleSingleReturn",
			v:    func(arg Int64) int64 { return 10 },
			args: []Value{Int64(10)},
			want: Int64(10),
		},
		{
			name: "SimpleMultiReturn",
			v:    func(arg Int64) (int64, string) { return 10, "hello" },
			args: []Value{Int64(10)},
			want: List{Items: []Value{Int64(10), String("hello")}},
		},
		{
			name: "PureVariadicNoCallArgs",
			v: func(args ...Int64) int64 {
				sum := int64(0)
				for _, arg := range args {
					sum += int64(arg)
				}
				return sum
			},
			want: Int64(0),
		},
		{
			name: "PureVariadicWithCallArgs",
			v: func(args ...Int64) int64 {
				sum := int64(0)
				for _, arg := range args {
					sum += int64(arg)
				}
				return sum
			},
			args: []Value{Int64(1), Int64(10)},
			want: Int64(11),
		},
		{
			name:    "ArityErrorNonVariadic",
			v:       func() {},
			args:    []Value{Int64(10)},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "ArityErrorWithVariadic",
			v:       func(first string, args ...int) {},
			args:    []Value{},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "ArgTypeMismatchNonVariadic",
			v:       func(a int) {},
			args:    []Value{String("hello")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "ArgTypeMismatchVariadic",
			v:       func(args ...int) {},
			args:    []Value{String("hello")},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			fn := reflectFn(reflect.ValueOf(tt.v))

			got, err := fn.Invoke(nil, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Invoke() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Invoke() got = %v, want %v", got, tt.want)
			}
		})
	}
}
