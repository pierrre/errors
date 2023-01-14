package errors

import (
	"github.com/pierrre/errors/errbase"
)

// New returns a new error with a message and a stack.
//
// Use fmt.Sprintf() to format the message.
func New(msg string) error {
	return newError(msg)
}

func newError(msg string) error {
	err := errbase.New(msg)
	err = stackSkip(err, 3)
	return err
}
