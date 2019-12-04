package main

import (
	"flag"
	"fmt"
	"log"
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
			log.Fatalf("failed to open file: %v", err)
		}
		defer fh.Close()

		result, err = sabre.ReadEval(scope, fh)
		if err != nil {
			log.Fatalf("failed to read-eval file content: %v", err)
		}
	}

	if executeStr != nil {
		result, err = sabre.ReadEvalStr(scope, *executeStr)
		if err != nil {
			log.Fatalf("failed to read-eval string: %v", err)
		}
	}

	fmt.Println(result)
}
