package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spy16/sabre"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Usage: sabre <src>")
		os.Exit(1)
	}

	mod, err := sabre.New("<arg>", strings.NewReader(os.Args[1])).All()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	for _, expr := range mod {
		fmt.Println(expr)
	}
}
