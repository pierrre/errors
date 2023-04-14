// Package errbase provides a simple error type.
package errbase

import (
	"fmt"
)

type base struct {
	msg string
}

// New creates a new error with the given message.
//
// It can be used as a "sentinel" error.
func New(msg string) error {
	return &base{
		msg: msg,
	}
}

// Newf creates a new error with the given formatted message.
//
// It doesn't support the %w verb.
func Newf(format string, args ...any) error {
	return New(fmt.Sprintf(format, args...))
}

func (err *base) Error() string {
	return err.msg
}
