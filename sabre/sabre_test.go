package sabre_test

import (
	"testing"

	"github.com/spy16/sabre/sabre"
)

func TestNew(t *testing.T) {
	rt := sabre.New()
	if rt == nil {
		t.Errorf("sabre.New() got nil")
	}
}
