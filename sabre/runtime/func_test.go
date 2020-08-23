package runtime

import (
	"regexp"
	"testing"
)

func TestGoFunc_Equals(t *testing.T) {
	f1 := GoFunc(func(env Runtime, args ...Value) (Value, error) {
		return Nil{}, nil
	})
	if !f1.Equals(f1) {
		t.Errorf("GoFunc.Equals() expecting true, got false")
	}

	f2 := GoFunc(nil)
	if f1.Equals(f2) {
		t.Errorf("GoFunc.Equals() expecting false, got true")
	}
}

func TestGoFunc_String(t *testing.T) {
	goFn := GoFunc(func(rt Runtime, args ...Value) (Value, error) {
		return nil, nil
	})

	pattern := regexp.MustCompile(`^GoFunc{0x[0-9a-f]+}$`)
	str := goFn.String()

	if !pattern.MatchString(str) {
		t.Errorf("GoFunc.String() expected result to match `%s`, got `%s`",
			pattern.String(), str)
	}
}
