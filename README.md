# Sabre

[![GoDoc](https://godoc.org/github.com/spy16/sabre?status.svg)](https://godoc.org/github.com/spy16/sabre) [![Go Report Card](https://goreportcard.com/badge/github.com/spy16/sabre)](https://goreportcard.com/report/github.com/spy16/sabre) [![Build Status](https://travis-ci.org/spy16/sabre.svg?branch=master)](https://travis-ci.org/spy16/sabre)

Sabre is highly customizable, embeddable LISP engine for Go.

## Features

* Highly Customizable reader/parser through a read table (Inspired by Clojure) (See [Reader](#reader))
* Built-in data types: nil, bool, string, number, character, keyword, symbol, list, vector, set, module
* Multiple number formats supported: decimal, octal, hexadecimal, radix and scientific notations.
* Full unicode support. Symbols can include unicode characters (Example: `find-δ`, `π` etc.)
* Character Literals with support for:
  1. simple literals  (e.g., `\a` for `a`)
  2. special literals (e.g., `\newline`, `\tab` etc.)
  3. unicode literals (e.g., `\u00A5` for `¥` etc.)
* Clojure style built-in special forms: `λ` or `fn*`, `def`, `if`, `do`, `throw`, `let*`
* Simple interface `sabre.Value` (and optional `sabre.Invokable`) for adding custom
  data types. (See [Evaluation](#evaluation))

## Usage

> Sabre requires Go 1.13 or higher.

### As Library

```go
package main

import "github.com/spy16/sabre"

func main() {
    scope := sabre.NewScope(nil)

    result, err := sabre.ReadEvalStr(scope, "(+ 1 2)")
    if err != nil {
        log.Fatalf("failed to eval: %v", err)
    }

    fmt.Printf("Result:\n %v\n", result)
}
```

See [Extending](#extending) for more information on customizing the reader or eval.

### Standalone

1. Install Sabre into `GOBIN` path: `go get -u -v github.com/spy16/sabre/cmd/sabre`
2. Run:
   1. `sabre` for REPL
   2. `sabre -e "(+ 1 2 3)"` for executing string
   3. `sabre -f "examples/full.lisp"` for executing file

> If you specify both `-f` and `-e` flags, file will be executed first and then the
> string will be executed in the same scope and you will be dropped into REPL. If
> REPL not needed, use `-norepl` option.

## Extending

### Reader

Sabre reader is inspired by Clojure reader and uses a _read table_. Reader supports
following forms:

* Numbers:
  * Integers use `int64` Go representation and can be specified using decimal, binary
    hexadecimal or radix notations. (e.g., 123, -123, 0b101011, 0xAF, 2r10100, 8r126 etc.)
  * Floating point numbers use `float64` Go representation and can be specified using
    decimal notation or scientific notation. (e.g.: 3.1412, -1.234, 1e-5, 2e3, 1.5e3 etc.)
* Characters: Characters use `rune` or `uint8` Go representation and can be written in 3 ways:
  * Simple: `\a`, `\λ`, `\β` etc.
  * Special: `\newline`, `\tab` etc.
  * Unicode: `\u1267`
* Boolean: `true` or `false` are converted to `Bool` type.
* Nil: `nil` is represented as a zero-allocation empty struct in Go.
* Keywords: Keywords are like symbols but start with `:` and evaluate to themselves.
* Symbols: Symbols can be used to name a value and can contain any Unicode symbol.
* Lists: Lists are zero or more forms contained within parenthesis. (e.g., `(1 2 3)`, `(1 [])`).
  Evaluating a list leads to an invocation.
* Vectors: Vectors are zero or more forms contained within brackets. (e.g., `[]`, `[1 2 3]`)
* Sets: Set is a container for zero or more unique forms. (e.g. `#{1 2 3}`)

Reader can be extended to add new syntactical features by adding _reader macros_
to the _read table_. _Reader Macros_ are implementations of `sabre.ReaderMacro`
function type. _Except numbers and symbols, everything else supported by the reader
is implemented using reader macros_.

### Evaluation

Eval logic for standard data types is fixed. But custom `sabre.Value` types can be
implemented to customize evaluation logic.

In addition, `sabre.Value` types can also implement `sabre.Invokable` interface to
enable invocation. For example `Vector` uses this to enable Clojure style element
access using `([1 2 3] 0)` (returns `1`)

> Please note that Sabre is _NOT_ an implementation of a particular LISP dialect (although
> it derives ideas from Clojure)

## TODO

* [x] Executor
* [x] Special Forms
* [X] REPL
* [ ] Standard Functions and Macros
* [ ] Optimizations
* [ ] Code Generation?
