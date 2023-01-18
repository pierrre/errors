package errors

import (
	"github.com/pierrre/errors/errmsg"
)

// Wrap adds a message to an error, and optionnally add a stack if it doesn't have one.
//
// See Message() and Stack() for more information.
func Wrap(err error, msg string) error {
	err = errmsg.Wrap(err, msg)
	err = ensureStack(err, 2)
	return err
}
