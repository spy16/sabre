package main

import (
	"fmt"
	"time"

	"github.com/spy16/sabre"
)

const rangeTco = `
(def range (fn* range [min max coll]
                 (if (= min max)
                   coll
                   (recur (inc min) max (coll.Conj min)))))

(print (range 0 10 '()))
(range 0 1000 '())
`

const rangeNotTco = `
(def range (fn* range [min max coll]
                 (if (= min max)
                   coll
                   (range (inc min) max (coll.Conj min)))))

(print (range 0 10 '()))
(range 0 1000 '())
`

func main() {
	scope := sabre.New()
	scope.BindGo("inc", inc)
	scope.BindGo("print", fmt.Println)
	scope.BindGo("=", sabre.Compare)

	initial := time.Now()
	_, err := sabre.ReadEvalStr(scope, rangeNotTco)
	if err != nil {
		panic(err)
	}
	final := time.Since(initial)
	fmt.Printf("no recur: %s\n", final)

	initial = time.Now()
	_, err = sabre.ReadEvalStr(scope, rangeTco)
	if err != nil {
		panic(err)
	}
	final = time.Since(initial)
	fmt.Printf("recur: %s\n", final)
}

func inc(val sabre.Int64) sabre.Int64 {
	return val + 1
}
