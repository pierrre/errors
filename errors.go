// Package errors provides error management.
//
// By convention, wrapping functions return a nil error if the given error is nil.
package errors

import (
	std_errors "errors"

	"github.com/pierrre/errors/errbase"
	"github.com/pierrre/errors/errmsg"
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

// Wrap adds a message to an error, and optionnally add a stack if it doesn't have one.
//
// See Message() and Stack() for more information.
func Wrap(err error, msg string) error {
	err = errstack.EnsureSkip(err, 1)
	err = errmsg.Wrap(err, msg)
	return err
}

// As calls std_errors.As.
//
// See https://pkg.go.dev/errors#As .
func As(err error, target any) bool {
	return std_errors.As(err, target)
}

// Is calls std_errors.Is.
//
// See https://pkg.go.dev/errors#Is .
func Is(err, target error) bool {
	return std_errors.Is(err, target)
}

// Unwrap calls std_errors.Unwrap.
//
// See https://pkg.go.dev/errors#Unwrap .
func Unwrap(err error) error {
	return std_errors.Unwrap(err)
}
