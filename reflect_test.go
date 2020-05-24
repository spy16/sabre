package sabre

import (
	"errors"
	"reflect"
	"testing"

	"github.com/spy16/sabre/core"
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
		want core.Value
	}{
		{
			name: "int64",
			v:    int64(10),
			want: core.Int64(10),
		},
		{
			name: "float",
			v:    float32(10.),
			want: core.Float64(10.),
		},
		{
			name: "uint8",
			v:    uint8('a'),
			want: core.Character('a'),
		},
		{
			name: "bool",
			v:    true,
			want: core.Bool(true),
		},
		{
			name: "Value",
			v:    core.Int64(10),
			want: core.Int64(10),
		},
		{
			name: "Nil",
			v:    nil,
			want: core.Nil{},
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
		getEnv  func() core.Env
		v       interface{}
		args    []core.Value
		want    core.Value
		wantErr bool
	}{
		{
			name: "WithEnvArgNoBinding",
			getEnv: func() core.Env {
				env := New()
				env.Bind("hello", core.Int64(10))
				return env
			},
			v: func(env core.Env) (core.Value, error) {
				return env.Eval(core.Symbol{Value: "hello"})
			},
			want:    core.Int64(10),
			wantErr: false,
		},
		{
			name: "SimpleNoArgNoReturn",
			v:    func() {},
			want: core.Nil{},
		},
		{
			name: "SimpleNoArg",
			v:    func() int { return 10 },
			want: core.Int64(10),
		},
		{
			name:    "NoArgSingleErrorReturn",
			v:       func() error { return errors.New("failed") },
			wantErr: true,
		},
		{
			name:    "NoArgSingleReturnNilError",
			v:       func() error { return nil },
			want:    core.Nil{},
			wantErr: false,
		},
		{
			name: "SimpleNoReturn",
			v:    func(arg core.Int64) {},
			args: []core.Value{core.Int64(10)},
			want: core.Nil{},
		},
		{
			name: "SimpleSingleReturn",
			v:    func(arg core.Int64) int64 { return 10 },
			args: []core.Value{core.Int64(10)},
			want: core.Int64(10),
		},
		{
			name: "MultiReturn",
			v: func(arg core.Int64) (int64, string) {
				return 10, "hello"
			},
			args: []core.Value{core.Int64(10)},
			want: core.Values([]core.Value{
				core.Int64(10),
				core.String("hello"),
			}),
		},
		{
			name:    "NoArgMultiReturnWithError",
			v:       func() (int, error) { return 0, errors.New("failed") },
			wantErr: true,
		},
		{
			name: "NoArgMultiReturnWithoutError",
			v:    func() (int, error) { return 10, nil },
			want: core.Int64(10),
		},
		{
			name: "PureVariadicNoCallArgs",
			v: func(args ...core.Int64) int64 {
				sum := int64(0)
				for _, arg := range args {
					sum += int64(arg)
				}
				return sum
			},
			want: core.Int64(0),
		},
		{
			name: "PureVariadicWithCallArgs",
			v: func(args ...core.Int64) int64 {
				sum := int64(0)
				for _, arg := range args {
					sum += int64(arg)
				}
				return sum
			},
			args: []core.Value{
				core.Int64(1),
				core.Int64(10),
			},
			want: core.Int64(11),
		},
		{
			name:    "ArityErrorNonVariadic",
			v:       func() {},
			args:    []core.Value{core.Int64(10)},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "ArityErrorWithVariadic",
			v:       func(first string, args ...int) {},
			args:    []core.Value{},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "ArgTypeMismatchNonVariadic",
			v:       func(a int) {},
			args:    []core.Value{core.String("hello")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "ArgTypeMismatchVariadic",
			v:       func(args ...int) {},
			args:    []core.Value{core.String("hello")},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			if tt.getEnv == nil {
				tt.getEnv = func() core.Env { return New() }
			}

			fn := reflectFn(reflect.ValueOf(tt.v))

			got, err := fn.Invoke(tt.getEnv(), tt.args...)
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
