// Package errjoin provides a Join function that combines multiple errors into a single error.
package errjoin

import (
	"github.com/pierrre/errors/errappend"
)

// Join returns an error that wraps the given errors.
// Any nil error values are discarded.
// Join returns nil if every value in errs is nil.
// The error formats as the concatenation of the strings obtained
// by calling the Error method of each element of errs, with a newline
// between each string.
//
// A non-nil error returned by Join implements the Unwrap() []error method.
// The errors may be inspected with [Is] and [As].
func Join(errs ...error) error {
	n := 0
	for _, err := range errs {
		if err != nil {
			n++
		}
	}
	if n == 0 {
		return nil
	}
	e := &joinError{
		errs: make([]error, 0, n),
	}
	for _, err := range errs {
		if err != nil {
			e.errs = append(e.errs, err)
		}
	}
	return e
}

type joinError struct {
	errs []error
}

func (e *joinError) Error() string {
	// Since Join returns nil if every value in errs is nil,
	// e.errs cannot be empty.
	if len(e.errs) == 1 {
		return e.errs[0].Error()
	}
	return errappend.String(e)
}

func (e *joinError) ErrorAppend(b []byte) []byte {
	for i, err := range e.errs {
		if i != 0 {
			b = append(b, '\n')
		}
		b = errappend.Append(b, err)
	}
	return b
}

func (e *joinError) Unwrap() []error {
	return e.errs
}
