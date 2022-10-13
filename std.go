package errors

import (
	std_errors "errors"
)

// This files contains aliases of the std errors package.
// So this package can be used as a drop-in replacement.

// As calls std_errors.As.
func As(err error, target any) bool {
	return std_errors.As(err, target)
}

// Is calls std_errors.Is.
func Is(err, target error) bool {
	return std_errors.Is(err, target)
}

// Unwrap calls std_errors.Unwrap.
func Unwrap(err error) error {
	return std_errors.Unwrap(err)
}
