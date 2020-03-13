# Sabre

[![GoDoc](https://godoc.org/github.com/spy16/sabre?status.svg)](https://godoc.org/github.com/spy16/sabre) [![Go Report Card](https://goreportcard.com/badge/github.com/spy16/sabre)](https://goreportcard.com/report/github.com/spy16/sabre) [![Build Status](https://travis-ci.org/spy16/sabre.svg?branch=master)](https://travis-ci.org/spy16/sabre)

Sabre is highly customizable, embeddable LISP engine for Go.

Check out [Slang](https://github.com/spy16/slang) for a tiny LISP written using *Sabre*.

## Features

* Highly Customizable reader/parser through a read table (Inspired by Clojure) (See [Reader](#reader))
* Built-in data types: nil, bool, string, number, character, keyword, symbol, list, vector, set,
  hash-map and module.
* Multiple number formats supported: decimal, octal, hexadecimal, radix and scientific notations.
* Full unicode support. Symbols can include unicode characters (Example: `find-δ`, `π` etc.)
  and `🧠`, `🏃` etc. (yes, smileys too).
* Character Literals with support for:
  1. simple literals  (e.g., `\a` for `a`)
  2. special literals (e.g., `\newline`, `\tab` etc.)
  3. unicode literals (e.g., `\u00A5` for `¥` etc.)
* Clojure style built-in special forms: `fn*`, `def`, `if`, `do`, `throw`, `let*`
* Simple interface `sabre.Value` and optional `sabre.Invokable`, `sabre.Seq` interfaces for
  adding custom data types. (See [Evaluation](#evaluation))
* A macro system.

> Please note that Sabre is _NOT_ an implementation of a particular LISP dialect. It provides
> pieces that can be used to build a LISP dialect or can be used as a scripting layer.

## Usage

What can you use it for?

1. Embedded script engine to provide dynamic behavior without requiring re-compilation
   of your application.
2. Business rule engine by exposing very specific & composable rule functions.
3. To build your own LISP dialect.

> Sabre requires Go 1.13 or higher.

### As Embedded Script Engine

Sabre has concept of `Scope` which is responsible for maintaining bindings. You can bind
any Go value and access it using LISP code, which makes it possible to expose parts of your
API and make it scriptable or build your own LISP dialect. Also, See [Extending](#extending)
for more information on customizing the reader or eval.

```go
package main

import "github.com/spy16/sabre"

func main() {
    scope := sabre.NewScope(nil)
    _ = scope.BindGo("inc", func(v int) int { return v+1 })

    result, _ := sabre.ReadEvalStr(scope, "(inc 10)")
    fmt.Printf("Result: %v\n", result) // should print "Result: 11"
}
```

### Expose through a REPL

Sabre comes with a tiny `repl` package that is very flexible and easy to setup
to expose your LISP through a read-eval-print-loop.

```go
package main

import (
  "context"

  "github.com/spy16/sabre"
  "github.com/spy16/sabre/repl"
)

func main() {
  scope := sabre.NewScope(nil)
  scope.BindGo("inc", func(v int) int { return v+1 })

  repl.New(scope,
    repl.WithBanner("Welcome to my own LISP!"),
    repl.WithPrompts("=>", "|"),
    // many more options available
  ).Loop(context.Background())
}
```

### Standalone

Sabre has a small reference LISP dialect named ***Slang*** (short for *Sabre Lang*) for
which a standalone binary is available. Check out [Slang](https://github.com/spy16/slang)
for instructions on installing *Slang*.

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
* HashMaps: HashMap is a container for key-value pairs (e.g., `{:name "Bob" :age 10}`)

Reader can be extended to add new syntactical features by adding _reader macros_
to the _read table_. _Reader Macros_ are implementations of `sabre.ReaderMacro`
function type. _Except numbers and symbols, everything else supported by the reader
is implemented using reader macros_.

### Evaluation

* `Keyword`, `String`, `Int`, `Float`, `Character`, `Bool`, `nil`, `MultiFn`,
  `Fn`, `Type` and `Any` evaluate to themselves.
* `Symbol` is resolved as follows:
  * If symbol has no `.`, symbol is directly used to lookup in current `Scope`
    to find the value.
  * If symbol is qualified (i.e., contains `.`), symbol is split using `.` as
    delimiter and first field is resolved as per previous rule and rest of the
    fields are recursively resolved as members. (For example, `foo.Bar.Baz`: `foo`
    is resolved from scope, `Bar` should be member of value of `foo`. And `Baz`
    should be member of value resolved for `foo.Bar`)
* Evaluating `HashMap`, `Vector` & `Set` simply yields new hashmap, vector and set
  whose values are evaluated values contained in the original hashmaap, vector and set.
* Evaluating `Module` evaluates all the forms in the module and returns the result
  of last evaluation. Any error stops the evaluation process.
* Empty `List` is returned as is.
* Non empty `List` is an invocation and evaluated using following rules:
  * If the first argument resolves to a special-form (`SpecialForm` Go type),
    it is invoked and return value is cached in the list. This return value
    is used for evaluating the list.
  * If the first argument resolves to a Macro, macro is invoked with the rest
    of the list as arguments and return value replaces the list with `(do retval)`
    form.
  * If first value resolves to an `Invokable` value, `Invoke()` is called. Functions
    are implemented using `MultiFn` which implements `Invokable`. `Vector` also implements
    `Invokable` and provides index access.
  * It is an error.
