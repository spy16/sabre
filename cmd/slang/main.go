package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/chzyer/readline"
	"github.com/spy16/sabre"
	"github.com/spy16/sabre/repl"
	"github.com/spy16/sabre/slang"
)

const help = `Slang %s [Commit: %s] [Compiled with %s]
Visit https://github.com/spy16/sabre for more.`

var (
	version = "N/A"
	commit  = "N/A"

	executeFile = flag.String("f", "", "File to read and execute")
	executeStr  = flag.String("e", "", "Execute string")
	noREPL      = flag.Bool("norepl", false, "Don't start REPL after executing file and string")
)

type temp struct {
	Name string
}

func (temp *temp) Foo() {
	fmt.Println("foo called")
}

func main() {
	flag.Parse()

	sl := slang.New()
	sl.BindGo("*version*", version)

	var result sabre.Value
	var err error

	if *executeFile != "" {
		fh, err := os.Open(*executeFile)
		if err != nil {
			fatalf("error: %v\n", err)
		}
		defer fh.Close()

		result, err = sl.ReadEval(fh)
		if err != nil {
			fatalf("error: %v\n", err)
		}
		sl.SwitchNS(sabre.Symbol{Value: "user"})
	}

	if *executeStr != "" {
		result, err = sl.ReadEvalStr(*executeStr)
		if err != nil {
			fatalf("error: %v\n", err)
		}
	}

	if *noREPL {
		fmt.Println(result)
		return
	}

	lr, errMapper := readlineInstance()

	repl := repl.New(sl,
		repl.WithBanner(fmt.Sprintf(help, version, commit, runtime.Version())),
		repl.WithInput(lr, errMapper),
		repl.WithOutput(lr.Stdout()),
		repl.WithPrompts("=>", "|"),
	)

	if err := repl.Loop(context.Background()); err != nil {
		fatalf("REPL exited with error: %v", err)
	}
	fmt.Println("Bye!")
}

func readlineInstance() (*readline.Instance, func(error) error) {
	lr, err := readline.New("")
	if err != nil {
		fatalf("readline: %v", err)
	}

	errMapper := func(e error) error {
		if errors.Is(e, readline.ErrInterrupt) {
			return nil
		}

		return e
	}

	return lr, errMapper
}

func fatalf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	os.Exit(1)
}
