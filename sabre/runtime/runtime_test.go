package runtime_test

import (
	"testing"

	"github.com/spy16/sabre/sabre/runtime"
)

func TestPosition_Pos_SetPos(t *testing.T) {
	p1 := runtime.Position{
		File:   "hello.lisp",
		Line:   1,
		Column: 100,
	}
	p2 := runtime.Position{}
	p2.SetPos(p1.Pos())

	if p1 != p2 {
		t.Errorf("expected p1 & p2 to be equal. p1=%+v, p2=%+v", p1, p2)
	}
}

func TestPosition_String(t *testing.T) {
	t.Parallel()

	t.Run("WithFile", func(t *testing.T) {
		p := runtime.Position{
			File:   "hello.lisp",
			Line:   1,
			Column: 100,
		}

		want := "hello.lisp:1:100"
		got := p.String()

		if got != want {
			t.Errorf("Position.String() want=`%s`, got=`%s`", want, got)
		}
	})

	t.Run("WithoutFile", func(t *testing.T) {
		p := runtime.Position{
			Line:   1,
			Column: 100,
		}

		want := "<unknown>:1:100"
		got := p.String()

		if got != want {
			t.Errorf("Position.String() want=`%s`, got=`%s`", want, got)
		}
	})
}
