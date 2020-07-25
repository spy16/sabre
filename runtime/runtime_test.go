package runtime_test

import (
	"errors"
	"testing"

	"github.com/spy16/sabre/runtime"
)

func TestNewErr(t *testing.T) {
	t.Parallel()

	t.Run("AlreadyError", func(t *testing.T) {
		cause := runtime.NewErr(false, runtime.Position{}, errors.New("failed"))

		pos := runtime.Position{
			File:   "hello.lisp",
			Line:   1,
			Column: 10,
		}
		err := runtime.NewErr(false, pos, cause)
		if err.Position != pos {
			t.Errorf("NewErr() expected position to be %+v, got=%+v", pos, err.Position)
		}
	})

	t.Run("AlreadyErrorPointer", func(t *testing.T) {
		cause := runtime.NewErr(false, runtime.Position{}, errors.New("failed"))

		pos := runtime.Position{
			File:   "hello.lisp",
			Line:   1,
			Column: 10,
		}
		err := runtime.NewErr(false, pos, &cause)
		if err.Position != pos {
			t.Errorf("NewErr() expected position to be %+v, got=%+v", pos, err.Position)
		}
	})
}

func TestError_Error(t *testing.T) {
	t.Parallel()

	t.Run("EvalError", func(t *testing.T) {
		err := runtime.NewErr(false, runtime.Position{}, errors.New("failed"))

		want := "eval-error in '<unknown>:0:0': failed"
		if err.Error() != want {
			t.Errorf("Error() want=`%s`, got=`%s`", want, err.Error())
		}
	})

	t.Run("ReadError", func(t *testing.T) {
		err := runtime.NewErr(true, runtime.Position{}, errors.New("failed"))

		want := "syntax error in '<unknown>:0:0': failed"
		if err.Error() != want {
			t.Errorf("Error() want=`%s`, got=`%s`", want, err.Error())
		}
	})
}
