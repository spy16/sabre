package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/spy16/sabre"
)

var executeFile = flag.String("f", "", "File to read and execute")
var executeStr = flag.String("e", "", "Execute string")

func main() {
	flag.Parse()

	scope := sabre.NewScope(nil)

	var result interface{}
	var err error

	if executeFile != nil && *executeFile != "" {
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

	if executeStr != nil {
		result, err = sabre.EvalStr(scope, *executeStr)
		if err != nil {
			fatalf("error: %v\n", err)
		}
	}

	fmt.Println(result)
}

func fatalf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	os.Exit(1)
}
