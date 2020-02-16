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

// NewREPL initializes a new Slang REPL and returns the instance.
func NewREPL(slang *Slang, opts ...REPLOption) (*REPL, error) {
	repl := REPL{
		sl:          slang,
		prompt:      "=>",
		multiPrompt: "|",
	}

	for _, option := range opts {
		option(&repl)
	}

	ri, err := readline.New(repl.prompt)
	if err != nil {
		return nil, err
	}
	repl.ri = ri

	return &repl, nil
}

// REPL implements a read-eval-print loop for Slang.
type REPL struct {
	sl          *Slang
	ri          *readline.Instance
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
	prompt := repl.prompt

	if multiline {
		nsPrefix = strings.Repeat(" ", len(nsPrefix)+1)
		prompt = repl.multiPrompt
	}

	repl.ri.SetPrompt(fmt.Sprintf("%s%s ", nsPrefix, prompt))
}

// REPLOption implmentations can be provided to NewREPL to configure the
// REPL during initialization.
type REPLOption func(repl *REPL)

// WithVimMode enables Vim based editing mode for REPL.
func WithVimMode() REPLOption {
	return func(repl *REPL) {
		repl.ri.SetVimMode(true)
	}
}

// WithBanner sets a welcome banner to print when the REPL starts.
func WithBanner(s string) REPLOption {
	return func(repl *REPL) {
		repl.banner = s
	}
}

// WithPrompts sets the prompt to be displayed when waiting for user input
// in the REPL.
func WithPrompts(oneLine, multiLine string) REPLOption {
	return func(repl *REPL) {
		repl.prompt = oneLine
		repl.multiPrompt = multiLine
	}
}
