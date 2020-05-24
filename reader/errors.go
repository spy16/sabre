package reader

import (
	"errors"
	"fmt"

	"github.com/spy16/sabre/core"
)

var (
	// ErrSkip can be returned by reader macro to indicate a no-op form which
	// should be discarded (e.g., Comments).
	ErrSkip = errors.New("skip expr")

	// ErrEOF is returned when stream ends prematurely to indicate that more
	// data is needed to complete the current form.
	ErrEOF = errors.New("unexpected EOF")
)

// Error wraps the parsing/eval errors with relevant information.
type Error struct {
	core.Position
	Cause  error
	Messag string
}

// Unwrap returns underlying cause of the error.
func (err Error) Unwrap() error {
	return err.Cause
}

func (err Error) Error() string {
	if e, ok := err.Cause.(Error); ok {
		return e.Error()
	}

	return fmt.Sprintf(
		"syntax error in '%s' (Line %d Col %d): %v",
		err.File, err.Line, err.Column, err.Cause,
	)
}
