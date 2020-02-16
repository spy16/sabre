package slang

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
	log "github.com/lthibault/log/pkg"
	"github.com/spy16/sabre"
)

const (
	promptPrefix    = "=>"
	multiLinePrompt = "->"
)

// NewREPL initializes a new Slang REPL and returns the instance.
func NewREPL(slang *Slang, opts ...REPLOption) *REPL {
	repl := REPL{
		runtime: slang,
	}

	for _, option := range withDefaults(opts) {
		option(&repl)
	}

	return &repl
}

// REPL implements a read-eval-print loop for Slang.
type REPL struct {
	Banner string

	log log.Logger

	runtime *Slang

	prompt Prompt
	// ri *readline.Instance

}

// Run starts the REPL loop and runs until the context is cancelled or
// a critical error occurs during ReadEval step.
func (repl *REPL) Run(ctx context.Context) (err error) {
	repl.prompt.SetPrompt(repl.getPrompt(promptPrefix))

	if repl.Banner != "" {
		fmt.Println(repl.Banner)
	}

	for {
		repl.prompt.SetPrompt(repl.getPrompt(promptPrefix))

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

			if form == nil {
				continue
			}

			repl.print(repl.runtime.Eval(form))
		}
	}
}

func (repl *REPL) print(res sabre.Value, err error) {
	if err != nil {
		fmt.Fprintf(os.Stdout, "error: %v\n", err)
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
			repl.prompt.SetPrompt(repl.getPrompt(multiLinePrompt))
		}

		line, err := repl.prompt.Readline()
		if err != nil {
			return nil, err
		}
		src += line + "\n"

		form, err = sabre.NewReader(strings.NewReader(src)).All()
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

func (repl *REPL) getPrompt(prompt string) string {
	return fmt.Sprintf("%s%s ", repl.runtime.CurrentNS(), prompt)
}
