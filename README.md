# Sabre

[![GoDoc](https://godoc.org/github.com/spy16/sabre?status.svg)](https://godoc.org/github.com/spy16/sabre) [![Go Report Card](https://goreportcard.com/badge/github.com/spy16/sabre)](https://goreportcard.com/report/github.com/spy16/sabre)

> WIP (See TODO)

Sabre is highly customizable, embeddable LISP engine for Go.

## Features

* Highly Customizable reader/parser through a read table (Inspired by Clojure)
* Built-in data types: string, number, character, keyword, symbol, list, vector
* Multiple number formats supported: decimal, octal, hexa-decimal, radix and scientific notations.
* Full unicode support.
  * Code can contain unicode characters.
  * Symbols can include unicode characters (Example: `find-δ`, `π` etc.)
* Character Literals with support for:
  1. simple literals  (e.g., `\a` for `a`)
  2. special literals (e.g., `\newline`, `\tab` etc.)
  3. unicode literals (e.g., `\u00A5` for `¥` etc.)

## Installation

To embed Sabre in your Go application to provide scripting capability, import `github.com/spy16/sabre`.

For standalone usage, run `go get -u -v github.com/spy16/sabre/cmd/sabre`

## TODO

* [ ] Executor
* [ ] Standard Functions, Special Forms and Macros
* [ ] REPL
* [ ] Optimizations
* [ ] Code Generation?

> Please note that Sabre is _NOT_ an implementation of a particular LISP dialect.
