// Package erriter provides a way to iterate over errors tree.
package erriter

// Iter iterates over an error tree.
//
// It explores the tree recursively, and calls f for each non-nil error.
func Iter(err error, f func(err error)) {
	for err != nil {
		f(err)
		var errs []error
		errs, err = Unwrap(err)
		for _, err := range errs {
			Iter(err, f)
		}
	}
}

// Unwrap unwraps an error.
//
// If the error implements `Unwrap() error`, it returns the unwrapped error.
// If the error implements `Unwrap() []error`, it returns the unwrapped errors.
// Otherwise, it returns nil.
func Unwrap(err error) ([]error, error) {
	switch err := err.(type) { //nolint:errorlint // We want to check which interface is implemented by the current error.
	case interface{ Unwrap() error }:
		return nil, err.Unwrap() //nolint:wrapcheck // We want to return the wrapped error.
	case interface{ Unwrap() []error }:
		return err.Unwrap(), nil
	}
	return nil, nil
}
