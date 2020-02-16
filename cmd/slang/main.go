package main

import (
	"context"

	"github.com/chzyer/readline"
	log "github.com/lthibault/log/pkg"
	"github.com/spy16/sabre/repl"
	"github.com/spy16/sabre/slang"
)

func main() {
	lr, err := readline.NewEx(&readline.Config{
		HistoryFile:       "/tmp/ww.history",
		HistorySearchFold: true,

		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})

	if err != nil {
		log.New().WithError(err).Fatal("readline initialization failed")
	}

	repl := repl.New(slang.New(),
		repl.WithPrompt(lr))

	if err := repl.Run(context.Background()); err != nil {
		log.New().WithError(err).Error("runtime error")
	}
}
