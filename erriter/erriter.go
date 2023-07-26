// Package erriter provides a way to iterate over errors tree.
package erriter

// Iter iterates over an error tree.
// It explores the tree recursively, and calls f for each error.
func Iter(err error, f func(err error)) {
	for err != nil {
		f(err)
		switch errw := err.(type) { //nolint:errorlint // We want to compare the current error.
		case interface{ Unwrap() error }:
			err = errw.Unwrap()
		case interface{ Unwrap() []error }:
			for _, err := range errw.Unwrap() {
				Iter(err, f)
			}
			err = nil
		default:
			err = nil
		}
	}
}
