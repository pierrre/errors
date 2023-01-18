// Package errtmp provides a way to mark errors as temporary.
package errtmp

import (
	"fmt"

	"github.com/pierrre/errors"
)

// Wrap marks an errors as temporary.
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

func (err *temporary) ErrorVerbose() string {
	return fmt.Sprintf("temporary = %t", err.tmp)
}

func (err *temporary) Temporary() bool {
	return err.tmp
}

// Is returns true if an error is temporary, false otherwise.
//
// By default, an error is temporary.
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
