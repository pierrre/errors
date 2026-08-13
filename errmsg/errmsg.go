// Package errmsg provides a way to add messages to errors.
package errmsg

import (
	"fmt"

	"github.com/pierrre/errors/errappend"
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

// Wrapf calls [Wrap] with a formatted message.
//
// It doesn't support the %w verb.
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
	return errappend.String(err)
}

func (err *message) ErrorAppend(b []byte) []byte {
	b = append(b, err.msg...)
	b = append(b, ": "...)
	b = errappend.Append(b, err.error)
	return b
}
