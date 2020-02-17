package repl

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// Option implementations can be provided to New() to configure the REPL
// during initialization.
type Option func(repl *REPL)

// WithInput sets the REPL's input stream. `nil` defaults to bufio.Scanner
// backed by os.Stdin
func WithInput(in Input, errMapper func(error) error) Option {
	if in == nil {
		in = &lineReader{
			scanner: bufio.NewScanner(os.Stdin),
			out:     os.Stdout,
		}
	}

	if errMapper == nil {
		errMapper = func(e error) error { return e }
	}

	return func(repl *REPL) {
		repl.input = in
		repl.inputErrMapper = errMapper
	}
}

// WithOutput sets the REPL's output stream.`nil` defaults to stdout.
func WithOutput(w io.Writer) Option {
	if w == nil {
		w = os.Stdout
	}

	return func(repl *REPL) {
		repl.output = w
	}
}

// WithBanner sets the REPL's banner which is displayed once when the REPL
// starts.
func WithBanner(banner string) Option {
	return func(repl *REPL) {
		repl.banner = strings.TrimSpace(banner)
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

func withDefaults(opts []Option) []Option {
	return append([]Option{
		WithInput(nil, nil),
		WithOutput(os.Stdout),
	}, opts...)
}

type lineReader struct {
	scanner *bufio.Scanner
	out     io.Writer
	prompt  string
}

func (lr *lineReader) Readline() (string, error) {
	lr.out.Write([]byte(lr.prompt))

	if !lr.scanner.Scan() {
		if lr.scanner.Err() == nil { // scanner swallows EOF
			return lr.scanner.Text(), io.EOF
		}

		return "", lr.scanner.Err()
	}

	return lr.scanner.Text(), nil
}

// no-op
func (lr *lineReader) SetPrompt(p string) {
	lr.prompt = p
}
