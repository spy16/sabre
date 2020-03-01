# Changelog

## v0.3.3 (2020-03-01)

* Prevent special forms being passed as value at `Symbol.Eval` level.
* Add checks for ambiguous arity detection in `MultiFn`.

## v0.3.2 (2020-02-28)

* Expand macro-expansions in place during 'List.parse()'
* Add `MacroExpand` public function
* Prevent macros being passed as value at `Symbole.Eval` level.
* All `Values` operations return `&List{}` as result instead of `Values`.

## v0.3.1 (2020-02-24)

* Move `Slang` to separate [repository](https://github.com/spy16/slang)
* Add support for member access using qualified symbols (i.e., `foo.Bar.Baz`).

## v0.3.0 (2020-02-22)

* Add support for macros through `macro*` special form.
* Use Macro support to add `defn` and `defmacro` macros.
* Rewrite `core.lisp` with `defn` macro.

## v0.2.3 (2020-02-22)

* Add support for custom special forms through `SpecialForm` type.
* Update package documentation.
* Remove all type specific functions in favour generic slang core.

## v0.2.2 (2020-02-21)

* Add type init through Type.Invoke() method
* Remove GoFunc in favor of Fn

## v0.2.1 (2020-02-19)

* Add slang tests using lisp files (#8)
* Added tests for function definitions
* Improve function call reflection logic to handle error returns

## v0.2.0 (2020-02-18)

* Add evaluation error with positional info
* Add position info to Set, List, Vector, Symbol
* Add slang runtime package, add generic repl package
* Add support for variadic functions

## v0.1.3 (2020-02-04)

* Add Values type and Seq types
* Add let and throw special forms
* Add support for multi-arity functions
* Convert List, Set, Vector, Symbol types to struct
* Modify List, Set, Vector types to embed Values type
* Move special form functions into sabre root package
* Add parent method to scope and modify def to apply at root scope

## v0.1.2 (2020-01-23)

* Add working clojure style quote system
* Move SpecialFn to sabre root package as GoFunc
* remove redundant strictFn type

## v0.1.1 (2020-01-20)

* Add error function and type functions
* Add experimental Set implementation
* Add special Nil type, add not & do core functions
* Add type check and type init functions for all types
* Add unit tests for all string and eval methods
* Fix nested lambda binding issue
* Split builtin functions into core package

## v0.1.0 (2020-01-18)

* Fully working LISP reader.
* Working Evaluation logic with no automatic type conversion.
* Core functions `def`, `eval` and `fn` implemented.
* Simple working REPL implemented.
