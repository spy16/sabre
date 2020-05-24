package sabre

import (
	"reflect"
	"testing"

	"github.com/spy16/sabre/core"
)

func TestReadEval(t *testing.T) {
	t.Parallel()

	table := []struct {
		title   string
		src     string
		getEnv  func() core.Env
		want    core.Value
		wantErr bool
	}{
		{
			title:   "ReadError",
			src:     "(hello",
			getEnv:  func() core.Env { return New() },
			want:    nil,
			wantErr: true,
		},
		{
			title: "Successful",
			src:   "pi",
			getEnv: func() core.Env {
				env := New()
				_ = env.Bind("pi", core.Float64(3.1412))
				return env
			},
			want:    core.Float64(3.1412),
			wantErr: false,
		},
	}

	for _, tt := range table {
		t.Run(tt.title, func(t *testing.T) {
			got, err := ReadEvalStr(tt.getEnv(), tt.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadEvalStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadEvalStr() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}
