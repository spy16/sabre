package runtime

import (
	"reflect"
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

func TestFn_Eval(t *testing.T) {
	fn := &Fn{
		Args:     []string{"a", "b"},
		Variadic: true,
		Body:     NewSeq(),
	}

	res, err := fn.Eval(nil)
	if err != nil {
		t.Errorf("Fn.Eval() unexpected error: %+v", err)
	}

	if !reflect.DeepEqual(fn, res) {
		t.Errorf("Fn.Eval() want=%+v, got=%+v", fn, res)
	}
}

func TestFn_Invoke(t *testing.T) {
	t.Parallel()
	table := []struct {
		title   string
		fn      *Fn
		getRT   func() Runtime
		args    []Value
		wantErr bool
		want    Value
	}{
		{
			title: "NoBody_NoArgs",
			fn:    &Fn{},
			want:  Nil{},
		},
	}

	for _, tt := range table {
		t.Run(tt.title, func(t *testing.T) {
			var rt Runtime
			if tt.getRT != nil {
				rt = tt.getRT()
			}

			got, err := tt.fn.Invoke(rt)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fn.Invoke() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fn.Invoke() got = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestFn_Equals(t *testing.T) {
	t.Parallel()

	fn := &Fn{
		Args:     []string{"a", "b"},
		Variadic: true,
		Body:     Float64(1.3),
	}

	table := []struct {
		title string
		other Value
		want  bool
	}{
		{
			title: "SameValue",
			other: fn,
			want:  true,
		},
		{
			title: "DifferentArgs",
			other: &Fn{
				Args:     []string{"b"},
				Variadic: true,
				Body:     Float64(1.3),
			},
			want: false,
		},
		{
			title: "NonVariadic",
			other: &Fn{
				Args:     fn.Args,
				Variadic: false,
				Body:     fn.Body,
			},
			want: false,
		},
		{
			title: "DifferentBody",
			other: &Fn{
				Args:     fn.Args,
				Variadic: fn.Variadic,
				Body:     String("something else"),
			},
			want: false,
		},
		{
			title: "NonFnValue",
			other: Float64(1.),
			want:  false,
		},
	}

	for _, tt := range table {
		t.Run(tt.title, func(t *testing.T) {
			got := fn.Equals(tt.other)
			if tt.want != got {
				t.Errorf("Fn.Equals() want=%+v, got=%+v", tt.want, got)
			}
		})
	}
}

func TestFn_String(t *testing.T) {
	t.Run("Variadic", func(t *testing.T) {
		fn := &Fn{
			Args:     []string{"a", "b"},
			Variadic: true,
			Body:     Float64(1.3),
		}

		want := "(fn [a & b] 1.300000)"
		got := fn.String()

		if want != got {
			t.Errorf("Fn.String() \nwant=`%s`\ngot =`%s`", want, got)
		}
	})

	t.Run("NonVariadic", func(t *testing.T) {
		fn := &Fn{
			Args:     []string{"a", "b"},
			Variadic: false,
			Body:     Float64(1.3),
		}

		want := "(fn [a b] 1.300000)"
		got := fn.String()

		if want != got {
			t.Errorf("Fn.String() \nwant=`%s`\ngot =`%s`", want, got)
		}
	})
}
