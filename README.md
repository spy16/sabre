# Sabre

[![GoDoc](https://godoc.org/github.com/spy16/sabre?status.svg)](https://godoc.org/github.com/spy16/sabre) [![Go Report Card](https://goreportcard.com/badge/github.com/spy16/sabre)](https://goreportcard.com/report/github.com/spy16/sabre) [![Build Status](https://travis-ci.org/spy16/sabre.svg?branch=master)](https://travis-ci.org/spy16/sabre)

> WIP (See TODO)

Sabre is highly customizable, embeddable LISP engine for Go.

## Features

* Highly Customizable reader/parser through a read table (Inspired by Clojure)
* Built-in data types: string, number, character, keyword, symbol, list, vector
* Multiple number formats supported: decimal, octal, hexadecimal, radix and scientific notations.
* Full unicode support. Symbols can include unicode characters (Example: `find-δ`, `π` etc.)
* Character Literals with support for:
  1. simple literals  (e.g., `\a` for `a`)
  2. special literals (e.g., `\newline`, `\tab` etc.)
  3. unicode literals (e.g., `\u00A5` for `¥` etc.)
* Simple evaluation logic with support for adding custom special-forms.

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
> string will be executed in the same scope.

## Extending

### Reader

Reader uses a macro table to allow adding support for new syntax. Reader macro is defined
as:

```go
type Macro func(rd *Reader, init rune) (Value, error)
```

For example following snippet adds support for Unix-like absolute file path.

```go
src := `/home/bob/documents/hello`

rd := sabre.NewReader(strings.NewReader(src))
rd.SetMacro('/', readUnixPath)
```

And the reader macro implementation:

```go
func readUnixPath(rd *sabre.Reader, init rune) (sabre.Value, error) {
    var path strings.Builder
    path.WriteRune(init)

    for {
        r, err := rd.NextRune()
        if err != nil {
            if err == io.EOF {
                break
            }

            return nil, err
        }

        if rd.IsTerminal(r) && r != '/' {
            break
        }

        path.WriteRune(r)
    }

    return sabre.String(path.String()), nil
}
```

### Evaluation

Eval logic for standard data types is fixed. But custom `sabre.Value` types can be
implemented to customize evaluation logic. Custom macros/special forms can be added
by wrapping a custom Go function with `SpecianFn` type.

> Please note that Sabre is _NOT_ an implementation of a particular LISP dialect.

## TODO

* [x] Executor
* [ ] Standard Functions, Special Forms and Macros
* [ ] REPL
* [ ] Optimizations
* [ ] Code Generation?
