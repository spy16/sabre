package main

import (
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/chzyer/readline"
	"github.com/spy16/sabre/repl"
	"github.com/spy16/sabre/slang"
)

const banner = "SLANG - a tiny lisp based on Sabre."

func main() {
	lr, err := readline.NewEx(&readline.Config{
		HistoryFile:       "/tmp/slang.history",
		HistorySearchFold: true,

		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})

	if err != nil {
		log.Fatal(err)
	}

	repl := repl.New(slang.New(),
		repl.WithBanner(banner),
		repl.WithInput(lr),
		repl.WithOutput(lr.Stdout()))

	if err := loop(repl); err != nil {
		log.Fatal(err)
	}
}

func loop(repl *repl.REPL) (err error) {
	for {
		if err = repl.Next(); err != nil {
			switch {
			case errors.Is(err, io.EOF):
				return nil
			case errors.Is(err, readline.ErrInterrupt):
				// print nothing
			default:
				fmt.Fprintf(repl, "%+v", err)
			}
		}
	}
}
