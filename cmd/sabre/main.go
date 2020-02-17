package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/chzyer/readline"
	"github.com/spy16/sabre/repl"
	"github.com/spy16/sabre/slang"
)

var (
	version = "N/A"
	commit  = "N/A"
)

const help = `Sabre %s [Commit: %s] [Compiled with %s]
Visit https://github.com/spy16/sabre for more.`

var executeFile = flag.String("f", "", "File to read and execute")
var executeStr = flag.String("e", "", "Execute string")
var noREPL = flag.Bool("norepl", false, "Don't start REPL after executing file and string")

func main() {
	flag.Parse()

	sl := slang.New()
	sl.BindGo("version", version)

	var result interface{}
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

	lr, err := readline.New("")
	if err != nil {
		fatalf("readline: %v", err)
	}

	repl := repl.New(sl,
		repl.WithBanner(fmt.Sprintf(help, version, commit, runtime.Version())),
		repl.WithInput(lr),
		repl.WithOutput(lr.Stdout()))

	if err := loop(repl); err != nil {
		fatalf("runtime error: %v", err)
	}
}

func fatalf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	os.Exit(1)
}

func loop(repl *repl.REPL) (err error) {
	for {
		if err = repl.Next(); err != nil {
			switch {
			case errors.Is(err, io.EOF):
				return nil
			case errors.Is(err, readline.ErrInterrupt):
				// continue
			default:
				fmt.Fprintf(repl, "%+v", err)
			}
		}
	}
}
