package sabre

import (
	"io"

	"github.com/spy16/sabre/sabre/collection"
	"github.com/spy16/sabre/sabre/core"
)

// NewReader returns an instance of sabre Reader for the given stream which
// can read primitive types and collection types.
func NewReader(r io.Reader) *core.Reader {
	macros := map[rune]core.ReaderMacro{
		'[': collection.VectorReader,
		']': core.UnmatchedDelimiter,
		'{': collection.HashMapReader,
		'}': core.UnmatchedDelimiter,
	}
	dispatchMacros := map[rune]core.ReaderMacro{}

	rd := core.NewReader(r)
	for init, macro := range macros {
		rd.SetMacro(init, macro, false)
	}

	for init, macro := range dispatchMacros {
		rd.SetMacro(init, macro, true)
	}
	return rd
}

// ReadEval reads from the given io reader using the default reader instance,
// evaluates all forms against the env and returns the result of the last eval.
func ReadEval(env core.Env, r io.Reader) (core.Value, error) {
	forms, err := NewReader(r).All()
	if err != nil {
		return nil, err
	}

	if len(forms) == 0 {
		return core.Nil{}, nil
	}

	res, err := core.EvalAll(env, forms)
	if err != nil {
		return nil, err
	}

	return res[len(res)-1], nil
}
