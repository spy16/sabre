// Package sabre builds a fully functional LISP environment using facilities
// provided by core package.
package sabre

import (
	"io"
	"strings"

	"github.com/spy16/sabre/core"
	"github.com/spy16/sabre/reader"
)

// New returns a new instance of root sabre environment.
func New() core.Env {
	env := core.New(nil)
	return env
}

// ReadEvalStr reads forms from src using a default reader instance and evaluates
// them against env.
func ReadEvalStr(env core.Env, src string) (core.Value, error) {
	rd := reader.New(strings.NewReader(src))
	form, err := rd.All()
	if err != nil {
		return nil, err
	}
	return env.Eval(form)
}

// ReadEval reads forms from r using a default reader instance and evaluates
// them against env.
func ReadEval(env core.Env, r io.Reader) (core.Value, error) {
	form, err := reader.New(r).All()
	if err != nil {
		return nil, err
	}
	return env.Eval(form)
}
