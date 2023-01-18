package errors

import (
	"github.com/pierrre/errors/errmsg"
	"github.com/pierrre/errors/errstack"
)

// Wrap adds a message to an error, and optionnally add a stack if it doesn't have one.
//
// See Message() and Stack() for more information.
func Wrap(err error, msg string) error {
	err = errstack.EnsureSkip(err, 1)
	err = errmsg.Wrap(err, msg)
	return err
}
