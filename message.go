package errors

import (
	"fmt"
)

// Message adds a message to an error.
//
// The error message is "<msg>: <err>".
//
// If the given message is empty, the returned error is the given error.
//
// Use fmt.Sprintf() to format the message.
func Message(err error, msg string) error {
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

type message struct {
	error
	msg string
}

func (err *message) Unwrap() error {
	return err.error
}

func (err *message) Error() string {
	return fmt.Sprintf("%s: %v", err.msg, err.error)
}
