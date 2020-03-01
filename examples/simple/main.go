package main

import (
	"fmt"

	"github.com/spy16/sabre"
)

const program = `
(def result (sum 1 2 3))
(printf "Sum of numbers is %s\n" result)
`

func main() {
	scope := sabre.New()
	scope.BindGo("sum", sum)
	scope.BindGo("printf", fmt.Printf)

	_, err := sabre.ReadEvalStr(scope, program)
	if err != nil {
		panic(err)
	}
}

func sum(nums ...float64) float64 {
	sum := 0.0
	for _, n := range nums {
		sum += n
	}

	return sum
}
