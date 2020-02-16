package repl

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/spy16/sabre"
)

// Runtime .
type Runtime interface {
	CurrentNS() string
	Eval(sabre.Value) (sabre.Value, error)
}

// New Read-Evaluate-Print Loop.
func New(r Runtime, opts ...Option) *REPL {
	repl := &REPL{runtime: r}

	for _, option := range withDefaults(opts) {
		option(repl)
	}

	return repl
}

// REPL implements a read-eval-print loop for Slang.
type REPL struct {
	runtime     Runtime
	input       Inputter
	banner      string
	prompt      string
	multiPrompt string
}

// Run starts the REPL loop and runs until the context is cancelled or
// a critical error occurs during ReadEval step.
func (repl *REPL) Run(ctx context.Context) (err error) {
	if repl.banner != "" {
		fmt.Println(repl.banner)
	}

	for {
		repl.setPrompt(false)

		select {
		case <-ctx.Done():
			return ctx.Err()

		default:
			form, err := repl.read()
			if err != nil {
				if err == readline.ErrInterrupt ||
					err == io.EOF {
					return nil
				}

				repl.print(nil, err)
				continue
			}

			if form != nil {
				repl.print(repl.runtime.Eval(form))
			}
		}
	}
}

func (repl *REPL) print(res sabre.Value, err error) {
	if err != nil {
		fmt.Fprintf(os.Stdout, "%v\n", err)
		return
	}

	fmt.Fprintf(os.Stdout, "%s\n", res)
}

func (repl *REPL) read() (sabre.Value, error) {
	var src string
	var form sabre.Value
	lineNo := 1

	for {
		if lineNo > 1 {
			repl.setPrompt(true)
		}

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

		form, err = rd.All()
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
