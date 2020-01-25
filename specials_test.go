package sabre_test

import (
	"testing"

	"github.com/spy16/sabre"
)

func TestLambda(t *testing.T) {
	t.Parallel()

	table := []struct {
		name    string
		args    []sabre.Value
		wantErr bool
	}{
		{
			name:    "InsufficientArgs",
			args:    nil,
			wantErr: true,
		},
		{
			name:    "InvalidArgList",
			args:    []sabre.Value{sabre.Int64(0), nil},
			wantErr: true,
		},
		{
			name: "NotSymbolVector",
			args: []sabre.Value{
				sabre.Vector{Items: []sabre.Value{sabre.Int64(1)}},
				sabre.Int64(10),
			},
			wantErr: true,
		},
		{
			name: "Successful",
			args: []sabre.Value{
				sabre.Vector{
					Items: []sabre.Value{sabre.Symbol{Value: "a"}, sabre.Symbol{Value: "b"}}},
				sabre.Int64(10),
			},
			wantErr: false,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sabre.Lambda(nil, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Lambda() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("Lambda() expecting non-nil, got nil")
				return
			}
		})
	}
}
