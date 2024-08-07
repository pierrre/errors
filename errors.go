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

// ErrUnsupported is an alias for [std_errors.ErrUnsupported].
var ErrUnsupported = std_errors.ErrUnsupported

// New returns a new error with a message and a stack.
//
// Warning: don't use this function to create a "sentinel" error, because it contains the stack of the (main) goroutine creating it.
// Use [errbase.New] instead.
func New(msg string) error {
	err := errbase.New(msg)
	err = errstack.WrapSkip(err, 1)
	return err //nolint: wrapcheck // The error is wrapped.
}

// Newf returns a new error with a formatted message and a stack.
//
// It supports the %w verb.
func Newf(format string, args ...any) error {
	err := errbase.Newf(format, args...)
	err = errstack.WrapSkip(err, 1)
	return err //nolint: wrapcheck // The error is wrapped.
}

// Wrap adds a message to an error, and a stack if it doesn't have one.
func Wrap(err error, msg string) error {
	err = errstack.EnsureSkip(err, 1)
	err = errmsg.Wrap(err, msg)
	return err
}

// Wrapf adds a formatted message to an error, and a stack if it doesn't have one.
//
// It doesn't support the %w verb.
func Wrapf(err error, format string, args ...any) error {
	err = errstack.EnsureSkip(err, 1)
	err = errmsg.Wrapf(err, format, args...)
	return err
}

// As is an alias for [std_errors.As].
func As(err error, target any) bool {
	return std_errors.As(err, target)
}

// Is is an alias for [std_errors.Is].
func Is(err, target error) bool {
	return std_errors.Is(err, target)
}

// Join calls [std_errors.Join] and adds a stack.
func Join(errs ...error) error {
	err := std_errors.Join(errs...)
	err = errstack.WrapSkip(err, 1)
	return err //nolint:wrapcheck // The error is wrapped.
}

// Unwrap is an alias for [std_errors.Unwrap].
func Unwrap(err error) error {
	return std_errors.Unwrap(err)
}
