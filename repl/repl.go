package repl

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/spy16/sabre"
)

// New returns a new instance of REPL with given runtime. Option values can
// be used to configure REPL input, output etc.
func New(r Runtime, opts ...Option) *REPL {
	repl := &REPL{runtime: r}
	for _, option := range withDefaults(opts) {
		option(repl)
	}
	return repl
}

// Input implementation is used by REPL to read user-input. See WithInput()
// REPL option to configure an Input.
type Input interface {
	SetPrompt(string)
	Readline() (string, error)
}

// Runtime implementation is used by REPL to evaluate user input.
type Runtime interface {
	// CurrentNS should return the current active name-space in the runtime.
	CurrentNS() string

	// Eval should evaluate the given form and return the result. Any error
	// returned by Eval will be formatted and printed to configured output.
	Eval(form sabre.Value) (sabre.Value, error)
}

// REPL implements a read-eval-print loop for a generic Runtime.
type REPL struct {
	runtime        Runtime
	input          Input
	inputErrMapper func(err error) error
	output         io.Writer

	banner      string
	prompt      string
	multiPrompt string
}

// Loop starts the read-eval-print loop. Loop runs until context is cancelled
// or input stream returns an irrecoverable error.
func (repl *REPL) Loop(ctx context.Context) error {
	repl.printBanner()
	repl.setPrompt(false)

	if repl.runtime == nil {
		return errors.New("runtime is not set")
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		default:
			err := repl.readEvalPrint()
			if err != nil {
				if err == io.EOF {
					return nil
				}

				return err
			}
		}
	}
}

// readEval reads one form from the input, evaluates it and prints the result.
func (repl *REPL) readEvalPrint() error {
	form, err := repl.read()
	if err != nil {
		return repl.inputErrMapper(err)
	}

	if form == nil {
		return nil
	}

	return repl.print(repl.runtime.Eval(form))
}

func (repl *REPL) Write(b []byte) (int, error) {
	return repl.output.Write(b)
}

func (repl *REPL) print(res sabre.Value, err error) error {
	if err != nil {
		_, err = fmt.Fprintf(repl, "%v\n", err)
		return err
	}

	_, err = fmt.Fprintf(repl, "%s\n", res)
	return err
}

func (repl *REPL) read() (sabre.Value, error) {
	var src string
	lineNo := 1

	for {
		repl.setPrompt(lineNo > 1)

		line, err := repl.input.Readline()
		if err != nil {
			return nil, err
		}
		src += line + "\n"

		if strings.TrimSpace(src) == "" {
			return nil, nil
		}

		rd := sabre.NewReader(strings.NewReader(src))
		rd.File = "REPL"

		form, err := rd.All()
		if err != nil {
			if errors.Is(err, sabre.ErrEOF) {
				lineNo++
				continue
			}

			return nil, err
		}

		return form, nil
	}
}

func (repl *REPL) setPrompt(multiline bool) {
	nsPrefix := repl.runtime.CurrentNS()
	prompt := repl.prompt

	if multiline {
		nsPrefix = strings.Repeat(" ", len(nsPrefix)+1)
		prompt = repl.multiPrompt
	}

	repl.input.SetPrompt(fmt.Sprintf("%s%s ", nsPrefix, prompt))
}

func (repl *REPL) printBanner() {
	if repl.banner != "" {
		fmt.Println(repl.banner)
	}
}
