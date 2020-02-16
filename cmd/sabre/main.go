package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"

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

	repl, err := slang.NewREPL(sl,
		slang.WithBanner(fmt.Sprintf(help, version, commit, runtime.Version())),
	)
	if err != nil {
		fatalf("failed to setup REPL: %v", err)
	}

	repl.Run(context.Background())
}

func fatalf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	os.Exit(1)
}
