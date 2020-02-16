package slang

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

const (
	promptPrefix    = "=>"
	multiLinePrompt = "|"
)

// NewREPL initializes a new Slang REPL and returns the instance.
func NewREPL(slang *Slang, opts ...REPLOption) *REPL {
	repl := REPL{
		sl: slang,
	}

	for _, option := range opts {
		option(&repl)
	}

	return &repl
}

// REPL implements a read-eval-print loop for Slang.
type REPL struct {
	sl     *Slang
	ri     *readline.Instance
	Banner string
}

// Run starts the REPL loop and runs until the context is cancelled or
// a critical error occurs during ReadEval step.
func (repl *REPL) Run(ctx context.Context) (err error) {
	repl.ri, err = readline.New(promptPrefix)
	if err != nil {
		return err
	}

	if repl.Banner != "" {
		fmt.Println(repl.Banner)
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
				repl.print(repl.sl.Eval(form))
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

		line, err := repl.ri.Readline()
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
	nsPrefix := repl.sl.CurrentNS()
	prompt := promptPrefix

	if multiline {
		nsPrefix = strings.Repeat(" ", len(nsPrefix)+1)
		prompt = multiLinePrompt
	}

	repl.ri.SetPrompt(fmt.Sprintf("%s%s ", nsPrefix, prompt))
}

// REPLOption implmentations can be provided to NewREPL to configure the
// REPL during initialization.
type REPLOption func(repl *REPL)
