package sabre_test

import (
	"reflect"
	"testing"

	"github.com/spy16/sabre"
)

var _ sabre.Scope = (*sabre.MapScope)(nil)

func TestMapScope_Get(t *testing.T) {
	tests := []struct {
		name     string
		symbol   string
		getScope func() *sabre.MapScope
		want     sabre.Value
		wantErr  bool
	}{
		{
			name:   "WithBinding",
			symbol: "hello",
			getScope: func() *sabre.MapScope {
				scope := sabre.NewScope(nil, false)
				_ = scope.Bind("hello", sabre.String("Hello World!"))
				return scope
			},
			want: sabre.String("Hello World!"),
		},
		{
			name:   "WithBindingInParent",
			symbol: "pi",
			getScope: func() *sabre.MapScope {
				parent := sabre.NewScope(nil, false)
				_ = parent.Bind("pi", sabre.Float64(3.1412))
				return sabre.NewScope(parent, false)
			},
			want: sabre.Float64(3.1412),
		},
		{
			name:   "WithNoBinding",
			symbol: "hello",
			getScope: func() *sabre.MapScope {
				return sabre.NewScope(nil, false)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := tt.getScope()

			got, err := scope.Resolve(tt.symbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("Resolve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resolve() got = %v, want %v", got, tt.want)
			}
		})
	}
}
