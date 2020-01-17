package sabre_test

import (
	"reflect"
	"testing"

	"github.com/spy16/sabre"
)

func TestEval(t *testing.T) {
	t.Parallel()

	table := []struct {
		name    string
		src     string
		want    sabre.Value
		wantErr bool
	}{
		{
			name: "Empty",
			src:  "",
			want: sabre.List(nil),
		},
		{
			name: "SingleForm",
			src:  "123",
			want: sabre.Int64(123),
		},
		{
			name: "MultiForm",
			src:  `123 [] ()`,
			want: sabre.List(nil),
		},
		{
			name:    "ReadError",
			src:     `123 [] (`,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sabre.EvalStr(nil, tt.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("EvalStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EvalStr() got = %v, want %v", got, tt.want)
			}
		})
	}
}
