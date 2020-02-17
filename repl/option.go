package repl

import (
	"bufio"
	"io"
	"os"
	"sync"
)

// Option implmentations can be provided to New to configure the
// REPL during initialization.
type Option func(repl *REPL)

// WithInput sets the REPL's input stream.  `nil` defaults to a bufio.Scanner backed
// by os.Stdin
func WithInput(in Input) Option {

	if in == nil {
		in = &lineReader{Reader: os.Stdin}
	}

	return func(repl *REPL) {
		repl.input = in
	}
}

// WithOutput sets the REPL's output stream.  `nil` defaults to os.Stdout.
func WithOutput(w io.Writer) Option {

	if w == nil {
		w = os.Stdout
	}

	return func(repl *REPL) {
		repl.output = w
	}
}

// WithBanner sets the REPL's banner.
func WithBanner(banner string) Option {
	return func(repl *REPL) {
		repl.banner = banner
	}
}

// WithPrompts sets the prompt to be displayed when waiting for user input
// in the REPL.
func WithPrompts(oneLine, multiLine string) Option {
	return func(repl *REPL) {
		repl.prompt = oneLine
		repl.multiPrompt = multiLine
	}
}

func withDefaults(opt []Option) []Option {
	return append([]Option{
		WithPrompts("=>", "|"),
		WithInput(nil),
		WithOutput(nil),
		// WithSomeOtherOption(...)
	}, opt...)
}

type lineReader struct {
	once    sync.Once
	scanner *bufio.Scanner
	io.Reader
}

func (lr *lineReader) Readline() (string, error) {
	lr.once.Do(func() {
		lr.scanner = bufio.NewScanner(lr.Reader)
	})

	if !lr.scanner.Scan() && lr.scanner.Err() == nil { // scanner swallows EOF
		return lr.scanner.Text(), io.EOF
	}

	return lr.scanner.Text(), nil
}

// no-op
func (lr *lineReader) SetPrompt(string) {}
