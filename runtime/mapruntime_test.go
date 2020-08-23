package runtime_test

import (
	"testing"

	"github.com/spy16/sabre/runtime"
)

func Test_mapEnv(t *testing.T) {
	parent := runtime.New(nil)
	_ = parent.Bind("π", runtime.Float64(3.1412))

	env := runtime.New(parent)
	_ = env.Bind("message", runtime.String("Hello World!"))

	t.Run("EvalNil", func(t *testing.T) {
		v, err := env.Eval(nil)
		if err != nil {
			t.Errorf("mapRuntime.resolve(\"message\"): unexpected error: %v", err)
		}
		want := runtime.Nil{}
		if !runtime.Equals(v, want) {
			t.Errorf("mapRuntime.resolve(\"message\") want=%+v, got=%+v", want, v)
		}
	})

	t.Run("resolve", func(t *testing.T) {
		v, err := env.Resolve("message")
		if err != nil {
			t.Errorf("mapRuntime.resolve(\"message\"): unexpected error: %v", err)
		}
		want := runtime.String("Hello World!")
		if !runtime.Equals(v, want) {
			t.Errorf("mapRuntime.resolve(\"message\") want=%+v, got=%+v", want, v)
		}
	})

	t.Run("ResolveFromParent", func(t *testing.T) {
		v, err := env.Resolve("π")
		if err != nil {
			t.Errorf("mapRuntime.resolve(\"π\"): unexpected error: %v", err)
		}
		want := runtime.Float64(3.1412)
		if !runtime.Equals(v, want) {
			t.Errorf("mapRuntime.resolve(\"π\") want=%+v, got=%+v", want, v)
		}
	})
}
