// Package errmsg provides a way to add messages to errors.
package errmsg

import (
	"fmt"
)

// Wrap adds a message to an error.
//
// The error message is "<msg>: <err>".
//
// If the given message is empty, the returned error is the given error.
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	if msg == "" {
		return err
	}
	return &message{
		error: err,
		msg:   msg,
	}
}

// Wrapf adds a formatted message to an error.
//
// See Wrap.
//
// it doesn't support the %w verb.
func Wrapf(err error, format string, args ...any) error {
	return Wrap(err, fmt.Sprintf(format, args...))
}

type message struct {
	error
	msg string
}

func (err *message) Unwrap() error {
	return err.error
}

func (err *message) Error() string {
	return err.msg + ": " + err.error.Error()
}
