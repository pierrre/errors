// Package errbase provides a simple error type.
package errbase

import (
	"fmt"

	std_errors "errors"
)

// New creates a new error with the given message.
//
// It can be used as a "sentinel" error.
func New(msg string) error {
	return std_errors.New(msg)
}

// Newf creates a new error with the given formatted message.
//
// It supports the %w verb.
func Newf(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}
