package sabre

import (
	"testing"
)

func Test_lambdaForm(t *testing.T) {
	t.Parallel()

	table := []struct {
		name    string
		args    []Value
		wantErr bool
	}{
		{
			name:    "InsufficientArgs",
			args:    nil,
			wantErr: true,
		},
		{
			name:    "InvalidArgList",
			args:    []Value{Int64(0), nil},
			wantErr: true,
		},
		{
			name: "NotSymbolVector",
			args: []Value{
				Vector{Values: []Value{Int64(1)}},
				Int64(10),
			},
			wantErr: true,
		},
		{
			name: "Successful",
			args: []Value{
				Vector{
					Values: []Value{Symbol{Value: "a"}, Symbol{Value: "b"}}},
				Int64(10),
			},
			wantErr: false,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got, err := lambdaForm(nil, tt.args)
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
