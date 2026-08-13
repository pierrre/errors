// Package errtmp provides a way to mark errors as temporary.
package errtmp

import (
	"strconv"

	"github.com/pierrre/errors"
	"github.com/pierrre/errors/errappend"
)

// Wrap marks an error as temporary.
//
// The verbose message is "temporary = <tmp>".
func Wrap(err error, tmp bool) error {
	if err == nil {
		return nil
	}
	return &temporary{
		error: err,
		tmp:   tmp,
	}
}

type temporary struct {
	error
	tmp bool
}

func (err *temporary) Unwrap() error {
	return err.error
}

func (err *temporary) ErrorAppend(b []byte) []byte {
	return errappend.Append(b, err.error)
}

func (err *temporary) ErrorVerboseAppend(b []byte) []byte {
	b = append(b, "temporary = "...)
	b = strconv.AppendBool(b, err.tmp)
	return b
}

func (err *temporary) Temporary() bool {
	return err.tmp
}

// Is returns true if an error is temporary, false otherwise.
//
// By default, an error is considered temporary.
// This is the opposite of the usual convention where an error that does not implement a Temporary() bool method is considered not temporary.
// To explicitly mark an error as not temporary, wrap it with [Wrap] and the value false.
func Is(err error) bool {
	var werr interface {
		Temporary() bool
	}
	ok := errors.As(err, &werr)
	if ok {
		return werr.Temporary()
	}
	return true
}
