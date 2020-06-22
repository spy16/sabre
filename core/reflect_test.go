package core

import (
	"errors"
	"reflect"
	"testing"

	"github.com/spy16/sabre/runtime"
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
		want runtime.Value
	}{
		{
			name: "int64",
			v:    int64(10),
			want: runtime.Int64(10),
		},
		{
			name: "float",
			v:    float32(10.),
			want: runtime.Float64(10.),
		},
		{
			name: "uint8",
			v:    uint8('a'),
			want: runtime.Char('a'),
		},
		{
			name: "bool",
			v:    true,
			want: runtime.Bool(true),
		},
		{
			name: "Value",
			v:    runtime.Int64(10),
			want: runtime.Int64(10),
		},
		{
			name: "Nil",
			v:    nil,
			want: runtime.Nil{},
		},
		{
			name: "ReflectType",
			v:    reflect.TypeOf(10),
			want: Type{T: reflect.TypeOf(10)},
		},
		{
			name: "Any",
			v:    anyVal,
			want: Any{V: anyValRV},
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
		getRT   func() runtime.Runtime
		v       interface{}
		args    []runtime.Value
		want    runtime.Value
		wantErr bool
	}{
		{
			name: "WithScopeArgNoBinding",
			getRT: func() runtime.Runtime {
				sc := runtime.New(nil)
				sc.Bind("hello", runtime.Int64(10))
				return sc
			},
			v: func(rt runtime.Runtime) (runtime.Value, error) {
				return rt.Resolve("hello")
			},
			want:    runtime.Int64(10),
			wantErr: false,
		},
		{
			name: "SimpleNoArgNoReturn",
			v:    func() {},
			want: runtime.Nil{},
		},
		{
			name: "SimpleNoArg",
			v:    func() int { return 10 },
			want: runtime.Int64(10),
		},
		{
			name:    "NoArgSingleErrorReturn",
			v:       func() error { return errors.New("failed") },
			wantErr: true,
		},
		{
			name:    "NoArgSingleReturnNilError",
			v:       func() error { return nil },
			want:    runtime.Nil{},
			wantErr: false,
		},
		{
			name: "SimpleNoReturn",
			v:    func(arg runtime.Int64) {},
			args: []runtime.Value{runtime.Int64(10)},
			want: runtime.Nil{},
		},
		{
			name: "SimpleSingleReturn",
			v:    func(arg runtime.Int64) int64 { return 10 },
			args: []runtime.Value{runtime.Int64(10)},
			want: runtime.Int64(10),
		},
		{
			name: "MultiReturn",
			v:    func(arg runtime.Int64) (int64, string) { return 10, "hello" },
			args: []runtime.Value{runtime.Int64(10)},
			want: runtime.NewSeq(runtime.Int64(10), runtime.String("hello")),
		},
		{
			name:    "NoArgMultiReturnWithError",
			v:       func() (int, error) { return 0, errors.New("failed") },
			wantErr: true,
		},
		{
			name: "NoArgMultiReturnWithoutError",
			v:    func() (int, error) { return 10, nil },
			want: runtime.Int64(10),
		},
		{
			name: "PureVariadicNoCallArgs",
			v: func(args ...runtime.Int64) int64 {
				sum := int64(0)
				for _, arg := range args {
					sum += int64(arg)
				}
				return sum
			},
			want: runtime.Int64(0),
		},
		{
			name: "PureVariadicWithCallArgs",
			v: func(args ...runtime.Int64) int64 {
				sum := int64(0)
				for _, arg := range args {
					sum += int64(arg)
				}
				return sum
			},
			args: []runtime.Value{runtime.Int64(1), runtime.Int64(10)},
			want: runtime.Int64(11),
		},
		{
			name:    "ArityErrorNonVariadic",
			v:       func() {},
			args:    []runtime.Value{runtime.Int64(10)},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "ArityErrorWithVariadic",
			v:       func(first string, args ...int) {},
			args:    []runtime.Value{},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "ArgTypeMismatchNonVariadic",
			v:       func(a int) {},
			args:    []runtime.Value{runtime.String("hello")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "ArgTypeMismatchVariadic",
			v:       func(args ...int) {},
			args:    []runtime.Value{runtime.String("hello")},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			if tt.getRT == nil {
				tt.getRT = func() runtime.Runtime { return runtime.New(nil) }
			}

			fn := reflectFn(reflect.ValueOf(tt.v))

			got, err := fn.Invoke(tt.getRT(), tt.args...)
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
