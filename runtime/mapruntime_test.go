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
			t.Errorf("mapRuntime.Resolve(\"message\"): unexpected error: %v", err)
		}
		want := runtime.Nil{}
		if !runtime.Equals(v, want) {
			t.Errorf("mapRuntime.Resolve(\"message\") want=%+v, got=%+v", want, v)
		}
	})

	t.Run("Resolve", func(t *testing.T) {
		v, err := env.Resolve("message")
		if err != nil {
			t.Errorf("mapRuntime.Resolve(\"message\"): unexpected error: %v", err)
		}
		want := runtime.String("Hello World!")
		if !runtime.Equals(v, want) {
			t.Errorf("mapRuntime.Resolve(\"message\") want=%+v, got=%+v", want, v)
		}
	})

	t.Run("ResolveFromParent", func(t *testing.T) {
		v, err := env.Resolve("π")
		if err != nil {
			t.Errorf("mapRuntime.Resolve(\"π\"): unexpected error: %v", err)
		}
		want := runtime.Float64(3.1412)
		if !runtime.Equals(v, want) {
			t.Errorf("mapRuntime.Resolve(\"π\") want=%+v, got=%+v", want, v)
		}
	})
}
