// Package repl provides a REPL implementation and options to expose Sabre
// features through a read-eval-print-loop.
package repl

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/spy16/sabre/core"
	"github.com/spy16/sabre/reader"
)

// New returns a new instance of REPL with given sabre Scope. Option values
// can be used to configure REPL input, output etc.
func New(env core.Env, opts ...Option) *REPL {
	repl := &REPL{
		env:   env,
		curNS: func() string { return "" },
	}

	if ns, ok := env.(NamespacedScope); ok {
		repl.curNS = ns.CurrentNS
	}

	for _, option := range withDefaults(opts) {
		option(repl)
	}

	return repl
}

// NamespacedScope can be implemented by sabre.Scope implementations to allow
// namespace based isolation (similar to Clojure). REPL will call CurrentNS()
// method to get the current Namespace and display it as part of input prompt.
type NamespacedScope interface {
	CurrentNS() string
}

// REPL implements a read-eval-print loop for a generic Runtime.
type REPL struct {
	env         core.Env
	input       Input
	output      io.Writer
	mapInputErr ErrMapper
	curNS       func() string
	factory     ReaderFactory

	banner      string
	prompt      string
	multiPrompt string

	printer func(io.Writer, interface{}) error
}

// Input implementation is used by REPL to read user-input. See WithInput()
// REPL option to configure an Input.
type Input interface {
	SetPrompt(string)
	Readline() (string, error)
}

// Loop starts the read-eval-print loop. Loop runs until context is cancelled
// or input stream returns an irrecoverable error (See WithInput()).
func (repl *REPL) Loop(ctx context.Context) error {
	repl.printBanner()
	repl.setPrompt(false)

	if repl.env == nil {
		return errors.New("scope is not set")
	}

	for ctx.Err() == nil {
		err := repl.readEvalPrint()
		if err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}
	}

	return ctx.Err()
}

// readEval reads one form from the input, evaluates it and prints the result.
func (repl *REPL) readEvalPrint() error {
	form, err := repl.read()
	if err != nil {
		switch err.(type) {
		case reader.Error, core.EvalError:
			repl.print(err)
		default:
			return err
		}
	}

	if form == nil {
		return nil
	}

	v, err := repl.env.Eval(form)
	if err != nil {
		return repl.print(err)
	}

	return repl.print(v)
}

func (repl *REPL) Write(b []byte) (int, error) {
	return repl.output.Write(b)
}

func (repl *REPL) print(v interface{}) error {
	return repl.printer(repl.output, v)
}

func (repl *REPL) read() (core.Value, error) {
	var src string
	lineNo := 1

	for {
		repl.setPrompt(lineNo > 1)

		line, err := repl.input.Readline()
		err = repl.mapInputErr(err)
		if err != nil {
			return nil, err
		}

		src += line + "\n"

		if strings.TrimSpace(src) == "" {
			return nil, nil
		}

		rd := repl.factory.NewReader(strings.NewReader(src))
		rd.File = "REPL"

		form, err := rd.All()
		if err != nil {
			if errors.Is(err, reader.ErrEOF) {
				lineNo++
				continue
			}

			return nil, err
		}

		return form, nil
	}
}

func (repl *REPL) setPrompt(multiline bool) {
	if repl.prompt == "" {
		return
	}

	nsPrefix := repl.curNS()
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
