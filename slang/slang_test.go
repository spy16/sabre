package slang_test

import (
	"testing"

	"github.com/spy16/sabre"
	"github.com/spy16/sabre/slang"
)

var _ sabre.Scope = (*slang.Slang)(nil)

func TestSlang_Bind(t *testing.T) {
	sl := slang.New()

	tests := []struct {
		name    string
		symbol  string
		ns      string
		wantErr bool
	}{
		{
			name:    "CrossNamespaceBindingValidation",
			symbol:  "core/not",
			ns:      "user",
			wantErr: true,
		},
		{
			name:    "BindingInCurrentNS",
			symbol:  "hello",
			ns:      "user",
			wantErr: false,
		},
		{
			name:    "UserBinding",
			symbol:  "user/hello",
			ns:      "user",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sl.Bind(tt.symbol, sabre.Nil{})
			if (err != nil) != tt.wantErr {
				t.Errorf("Bind() error = %#v, wantErr %#v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestSlang_Resolve(t *testing.T) {
	sl := slang.New()
	sl.Bind("pi", sabre.Float64(3.1412))

	tests := []struct {
		name    string
		symbol  string
		wantErr bool
	}{
		{
			name:   "CoreBinding",
			symbol: "core/not",
		},
		{
			name:    "UserBinding",
			symbol:  "hello",
			wantErr: true,
		},
		{
			name:    "MissingUserBinding",
			symbol:  "hello",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := sl.Resolve(tt.symbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("Resolve() error = %#v, wantErr %#v", err, tt.wantErr)
				return
			}
		})
	}
}
