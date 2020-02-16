package repl

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/spy16/sabre"
)

// Input controls input to a REPL
type Input interface {
	SetPrompt(string)
	Readline() (string, error)
}

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
	once sync.Once

	runtime Runtime
	input   Input
	output  io.Writer

	banner      string
	prompt      string
	multiPrompt string
}

// Next form runs through once read-eval-print cycle, returning any errors encountered.
// It is safe to call Next() again after an error, unless the error is EOF.
func (repl *REPL) Next() error {
	repl.once.Do(func() {
		fmt.Println(repl.banner)
	})

	repl.setPrompt(false)

	form, err := repl.read()
	if err != nil {
		return err
	}

	if form != nil {
		return repl.print(repl.runtime.Eval(form))
	}

	return nil
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

func (repl *REPL) Write(b []byte) (int, error) {
	return repl.output.Write(b)
}
