// Package erriter provides a way to iterate over errors tree.
package erriter

// Iter iterates over an error tree.
// It explores the tree recursively, and calls f for each error.
func Iter(err error, f func(err error)) {
	for err != nil {
		f(err)
		err = Unwrap(err, func(errs []error) {
			for _, err := range errs {
				Iter(err, f)
			}
		})
	}
}

// Unwrap unwraps an error.
//
// If the error implements `Unwrap() error`, it returns the unwrapped error.
// If the error implements `Unwrap() []error`, it calls onErrs with the unwrapped errors, and returns nil.
// Otherwise, it returns nil.
func Unwrap(err error, onErrs func(errs []error)) error {
	switch err := err.(type) { //nolint:errorlint // We want to compare the current error.
	case interface{ Unwrap() error }:
		return err.Unwrap() //nolint:wrapcheck // We want to return the wrapped error.
	case interface{ Unwrap() []error }:
		onErrs(err.Unwrap())
	}
	return nil
}
