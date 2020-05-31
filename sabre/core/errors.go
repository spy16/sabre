package core

import (
	"errors"
	"fmt"
)

var (
	// ErrSkip can be returned by reader macro to indicate a no-op form which
	// should be discarded (e.g., Comments).
	ErrSkip = errors.New("skip expr")

	// ErrEOF is returned by reader when stream ends prematurely to indicate
	// that more data is needed to complete the current form.
	ErrEOF = errors.New("unexpected EOF")

	// ErrNotFound should be returned by an env implementation when a binding is
	// not found or by values that implement Associative when an entry is not
	// found.
	ErrNotFound = errors.New("not found")
)

// NewErr returns a sabre error object with given err as cause. If err is already
// a sabre Error, simply returns copy of it with given position attached.
func NewErr(isRead bool, pos Position, err error) Error {
	if ee, ok := err.(Error); ok {
		ee.Position = pos
		return ee
	} else if ee, ok := err.(*Error); ok && ee != nil {
		err := *ee
		err.Position = pos
		return err
	}

	return Error{
		Position:  pos,
		Cause:     err,
		IsReadErr: isRead,
	}
}

// Error represents errors during read or evaluation stages.
type Error struct {
	Position
	IsReadErr bool
	Message   string
	Cause     error
	Form      Value
}

// Unwrap returns the underlying cause of this error.
func (err Error) Unwrap() error { return err.Cause }

func (err Error) Error() string {
	if e, ok := err.Cause.(Error); ok {
		return e.Error()
	}

	if err.IsReadErr {
		return fmt.Sprintf(
			"syntax error in '%s' (Line %d Col %d): %v",
			err.File, err.Line, err.Column, err.Cause,
		)
	}

	return fmt.Sprintf("eval-error in '%s' (at line %d:%d): %v",
		err.File, err.Line, err.Column, err.Cause,
	)
}
