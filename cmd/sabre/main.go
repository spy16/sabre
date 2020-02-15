package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/spy16/sabre/slang"
)

var (
	version = "N/A"
	commit  = "N/A"
)

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

	slang.NewREPL(sl).Run(context.Background())
}

func fatalf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	os.Exit(1)
}
