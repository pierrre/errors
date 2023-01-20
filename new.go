package errors

import (
	"github.com/pierrre/errors/errbase"
	"github.com/pierrre/errors/errstack"
)

// New returns a new error with a message and a stack.
//
// Use fmt.Sprintf() to format the message.
func New(msg string) error {
	err := errbase.New(msg)
	err = errstack.WrapSkip(err, 1)
	return err //nolint: wrapcheck // The error is wrapped.
}
