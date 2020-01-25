package core_test

import (
	"reflect"
	"testing"

	"github.com/spy16/sabre"
	"github.com/spy16/sabre/core"
)

func TestCore(t *testing.T) {
	t.Parallel()

	table := []struct {
		name     string
		fn       sabre.Invokable
		args     []sabre.Value
		getScope func() sabre.Scope
		want     sabre.Value
		wantErr  bool
	}{
		{
			name:    "Not_InsufficientArgs",
			fn:      core.Fn(core.Not),
			args:    []sabre.Value{},
			wantErr: true,
		},
		{
			name: "Not_Nil",
			fn:   core.Fn(core.Not),
			args: []sabre.Value{sabre.Nil{}},
			want: sabre.Bool(true),
		},
		{
			name: "Not_Integer",
			fn:   core.Fn(core.Not),
			args: []sabre.Value{sabre.Int64(10)},
			want: sabre.Bool(false),
		},
		{
			name: "Not_False",
			fn:   core.Fn(core.Not),
			args: []sabre.Value{sabre.Bool(false)},
			want: sabre.Bool(true),
		},
		{
			name: "Not_True",
			fn:   core.Fn(core.Not),
			args: []sabre.Value{sabre.Bool(true)},
			want: sabre.Bool(false),
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			var scope sabre.Scope
			if tt.getScope != nil {
				scope = tt.getScope()
			}

			got, err := tt.fn.Invoke(scope, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLambdaFn(t *testing.T) {
	fn := sabre.LambdaFn(nil,
		[]sabre.Symbol{sabre.Symbol{Value: "arg1"}},
		[]sabre.Value{sabre.Symbol{Value: "arg1"}},
	)

	arg1Val := sabre.Int64(10)

	got, err := fn.Invoke(nil, arg1Val)
	if err != nil {
		t.Errorf("Invoke() unexpected error: %v", err)
	}

	if !reflect.DeepEqual(got, arg1Val) {
		t.Errorf("Invoke() want=%#v, got=%#v", arg1Val, got)
	}
}
