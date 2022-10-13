package errors

import (
	"fmt"
)

// New returns a new error with a message and a stack.
func New(msg string) error {
	return newError(msg)
}

// Newf returns a new error with a formatted message and a stack.
func Newf(format string, args ...any) error {
	msg := fmt.Sprintf(format, args...)
	return newError(msg)
}

func newError(msg string) error {
	err := newBase(msg)
	err = stackSkip(err, 3)
	return err
}
