// Package erriter provides a way to iterate over errors tree.
package erriter

import (
	"iter"
)

// All returns an [iter.Seq] that iterates over an error tree recursively.
func All(err error) iter.Seq[error] {
	return func(yield func(error) bool) {
		iterFunc(err, func(err error) bool {
			return yield(err)
		})
	}
}

func iterFunc(err error, f func(err error) bool) bool {
	for err != nil {
		ok := f(err)
		if !ok {
			return false
		}
		var errs []error
		errs, err = Unwrap(err)
		for _, err := range errs {
			ok := iterFunc(err, f)
			if !ok {
				return false
			}
		}
	}
	return true
}

// Unwrap unwraps an error.
//
// If the error implements `Unwrap() error`, it returns the unwrapped error.
// If the error implements `Unwrap() []error`, it returns the unwrapped errors.
// Otherwise, it returns nil.
func Unwrap(err error) ([]error, error) {
	switch err := err.(type) { //nolint:errorlint // We want to check which interface is implemented by the current error.
	case interface{ Unwrap() error }:
		return nil, err.Unwrap()
	case interface{ Unwrap() []error }:
		return err.Unwrap(), nil
	}
	return nil, nil
}
