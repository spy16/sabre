package runtime

import "fmt"

// NewErr returns a sabre error object with given err as cause. If err is already
// a sabre Error, simply returns copy of it with given position attached.
func NewErr(isRead bool, pos Position, err error) Error {
	if ee, ok := toErr(err); ok {
		ee.IsReadErr = isRead
		ee.Position = pos
		return ee
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
	if err.File == "" {
		err.File = "<unknown>"
	}

	if err.IsReadErr {
		return fmt.Sprintf("syntax error in '%s': %v", err.Position, err.Cause)
	}

	return fmt.Sprintf("eval-error in '%s': %v", err.Position, err.Cause)
}

func toErr(err error) (Error, bool) {
	switch e := err.(type) {
	case Error:
		return e, true

	case *Error:
		return *e, true

	default:
		return Error{}, false
	}
}
