package main

import (
	"context"
	"fmt"

	"github.com/spy16/sabre"
	"github.com/spy16/sabre/repl"
)

const program = `
(def result (sum 1 2 3))
(printf "Sum of numbers is %s\n" result)
`

func main() {
	rt := sabre.New()
	_ = rt.Bind("sum", sabre.ValueOf(sum))
	_ = rt.Bind("printf", sabre.ValueOf(fmt.Printf))

	_ = repl.New(rt,
		repl.WithPrompts("=>", ">"),
	).Loop(context.Background())
}

func sum(nums ...float64) float64 {
	sum := 0.0
	for _, n := range nums {
		sum += n
	}

	return sum
}
