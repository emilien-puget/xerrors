package xerrors

import (
	std_errors "errors"
)

func newErrorString(text string) error {
	return &errorString{text}
}

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

// New returns a newErrorString error with a message and a stack.
func New(msg string) error {
	err := newErrorString(msg)
	err = withStack(err, 2)
	return err
}

// As calls std_errors.As.
func As(err error, target any) bool {
	return std_errors.As(err, target)
}

// Is calls std_errors.Is.
func Is(err, target error) bool {
	return std_errors.Is(err, target)
}

// Unwrap calls std_errors.Unwrap.
func Unwrap(err error) error {
	return std_errors.Unwrap(err)
}
