package main

import (
	"context"

	log "github.com/lthibault/log/pkg"
	"github.com/spy16/sabre/slang"
)

func main() {
	lr, err := rl.NewEx(&rl.Config{
		HistoryFile:       "/tmp/ww.history",
		HistorySearchFold: true,

		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})

	if err != nil {
		log.New().WithError(err).Fatal("readline initialization failed")
	}

	repl := slang.NewREPL(slang.New(),
		slang.WithPrompt(lr))

	if err := repl.Run(context.Background()); err != nil {
		log.New().WithError(err).Error("runtime error")
	}
}
