package runtime_test

import (
	"reflect"
	"testing"

	"github.com/spy16/sabre/runtime"
)

func TestList_Eval(t *testing.T) {
	t.Parallel()

	runEvalTests(t, []evalTestCase{
		{
			title: "EmptyList",
			form:  runtime.NewSeq(),
			want:  runtime.NewSeq(),
		},
		{
			title:   "FirstEvalFailure",
			getEnv:  func() runtime.Runtime { return runtime.New(nil) },
			form:    runtime.NewSeq(runtime.Symbol{Value: "non-existent"}),
			wantErr: true,
		},
		{
			title:   "NonInvokable",
			getEnv:  func() runtime.Runtime { return runtime.New(nil) },
			form:    runtime.NewSeq(runtime.Int64(0)),
			wantErr: true,
		},
		{
			title:  "InvokableNoArgs",
			getEnv: func() runtime.Runtime { return runtime.New(nil) },
			form: runtime.NewSeq(runtime.GoFunc(func(env runtime.Runtime, args ...runtime.Value) (runtime.Value, error) {
				return runtime.String("called"), nil
			})),
			want:    runtime.String("called"),
			wantErr: false,
		},
		{
			title:  "InvokableWithArgs",
			getEnv: func() runtime.Runtime { return runtime.New(nil) },
			form: runtime.NewSeq(runtime.GoFunc(func(env runtime.Runtime, args ...runtime.Value) (runtime.Value, error) {
				return args[0], nil
			}), runtime.String("hello")),
			want:    runtime.String("hello"),
			wantErr: false,
		},
	})
}

func Test_linkedList(t *testing.T) {
	t.Parallel()

	seq1 := runtime.NewSeq()
	assert(t, seq1.Count() == 0, "Seq.Count() expected 0, got %d", seq1.Count())

	conjCount := seq1.Conj(runtime.Nil{}).Count()
	assert(t, conjCount == 1, "Seq.Count() expected 1, got %d", conjCount)

	seq2 := seq1.Cons(runtime.Int64(0))
	assert(t, seq2.Count() == 1, "Seq.Count() expected 1, got %d", seq2.Count())
	assert(t, seq1.Count() == 0, "Seq.Count() expected 0, got %d", seq1.Count())

	got, want := seq2.First(), runtime.Value(runtime.Int64(0))
	assert(t, runtime.Equals(want, got), "Seq.First() want=%+v, got=%+v", want, got)

	got, want = seq1.First(), runtime.Value(runtime.Nil{})
	assert(t, runtime.Equals(want, got), "Seq.First() want=%+v, got=%+v", want, got)

	seq3 := seq2.Conj(runtime.Keyword("hello"), runtime.String("foo"))
	assert(t, seq3.Count() == 3, "Seq.Count() expected 3, got %d", seq3.Count())
}

func TestList_String(t *testing.T) {
	l := runtime.NewSeq(
		runtime.Bool(true),
		runtime.Int64(10),
		runtime.Float64(3.1412),
	)

	want := `(true 10 3.141200)`
	got := l.String()

	if want != got {
		t.Errorf("List.String() want=`%s`, got=`%s`",
			want, got)
	}
}

func assert(t *testing.T, cond bool, msg string, args ...interface{}) {
	if !cond {
		t.Errorf(msg, args...)
	}
}

type evalTestCase struct {
	title   string
	getEnv  func() runtime.Runtime
	form    runtime.Value
	want    runtime.Value
	wantErr bool
}

func runEvalTests(t *testing.T, cases []evalTestCase) {
	for _, tt := range cases {
		t.Run(tt.title, func(t *testing.T) {
			var env runtime.Runtime
			if tt.getEnv != nil {
				env = tt.getEnv()
			}

			got, err := tt.form.Eval(env)
			if (err != nil) != tt.wantErr {
				t.Errorf("%s.Eval() error = %v, wantErr %v",
					reflect.TypeOf(tt.form), err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s.Eval() got = %v, want %v",
					reflect.TypeOf(tt.form), got, tt.want)
			}
		})
	}
}
