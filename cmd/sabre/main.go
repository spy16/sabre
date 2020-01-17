package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/spy16/sabre"
)

var executeFile = flag.String("f", "", "File to read and execute")
var executeStr = flag.String("e", "", "Execute string")
var startREPL = flag.Bool("repl", false, "Start REPL after executing file and string")

func main() {
	flag.Parse()

	scope := sabre.NewScope(nil, true)
	scope.Bind("version", sabre.String(version))

	var result interface{}
	var err error

	if *executeFile != "" {
		fh, err := os.Open(*executeFile)
		if err != nil {
			fatalf("error: %v\n", err)
		}
		defer fh.Close()

		result, err = sabre.Eval(scope, fh)
		if err != nil {
			fatalf("error: %v\n", err)
		}
	}

	if *executeStr != "" {
		result, err = sabre.EvalStr(scope, *executeStr)
		if err != nil {
			fatalf("error: %v\n", err)
		}
	}

	if *startREPL || (*executeFile == "" && *executeStr == "") {
		repl, err := newREPL(scope)
		if err != nil {
			fatalf("REPL: %v", err)
		}

		repl.Start(context.Background())
	} else {
		fmt.Println(result)
	}
}

func fatalf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	os.Exit(1)
}
