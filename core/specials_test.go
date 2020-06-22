package core_test

import (
	"reflect"
	"testing"

	"github.com/spy16/sabre/core"
	"github.com/spy16/sabre/runtime"
)

func TestDef(t *testing.T) {
	t.Parallel()

	t.Run("InvalidArity", func(t *testing.T) {
		_, err := core.Def(runtime.New(nil))
		if err == nil {
			t.Errorf("expecting an error, got nil")
		}
	})

	t.Run("NotSymbol", func(t *testing.T) {
		_, err := core.Def(runtime.New(nil), runtime.Int64(0), runtime.String("value"))
		if err == nil {
			t.Errorf("expecting an error, got nil")
		}
	})

	t.Run("ExprEvalFailed", func(t *testing.T) {
		_, err := core.Def(runtime.New(nil), runtime.Symbol{Value: "foo"}, runtime.Symbol{Value: "test"})
		if err == nil {
			t.Errorf("expecting an error, got nil")
		}
	})

	t.Run("Successful", func(t *testing.T) {
		val, err := core.Def(runtime.New(runtime.New(nil)),
			runtime.Symbol{Value: "foo"}, runtime.String("value"))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if !runtime.Equals(val, runtime.Symbol{Value: "foo"}) {
			t.Errorf("expected return value to be symbol 'foo', got %s", reflect.TypeOf(val))
		}
	})
}
