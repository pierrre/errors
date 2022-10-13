// Package errignore provides a way to mark errors as ignored.
package errignore

import (
	"fmt"
	"io"

	"github.com/pierrre/errors"
)

// Wrap marks an error as ignored.
//
// The verbose message is "ignored".
func Wrap(err error) error {
	if err == nil {
		return nil
	}
	return &ignore{
		error: err,
	}
}

type ignore struct {
	error
}

func (err *ignore) Unwrap() error {
	return err.error
}

func (err *ignore) ErrorVerbose(w io.Writer) {
	_, _ = fmt.Fprint(w, "ignored\n")
}

func (err *ignore) Ignored() bool {
	return true
}

// Is returns true if an error is ignored, false otherwise.
//
// By default, an error is not ignored.
func Is(err error) bool {
	var werr *ignore
	ok := errors.As(err, &werr)
	if ok {
		return werr.Ignored()
	}
	return false
}
