// Package errval provides a way to add values to errors.
package errval

import (
	"fmt"

	"github.com/pierrre/errors/erriter"
)

// VerboseStringer returns the string representation of a value in a verbose message.
//
// It can be changed in order to customize how values are formatted.
var VerboseStringer = func(v any) string {
	return fmt.Sprint(v)
}

// Wrap adds a value to an error.
//
// The verbose message is "value <key> = <val>".
// The value is written using the VerboseWriter function.
func Wrap(err error, key string, val any) error {
	if err == nil {
		return nil
	}
	return &value{
		error: err,
		key:   key,
		val:   val,
	}
}

type value struct {
	error
	key string
	val any
}

func (err *value) Unwrap() error {
	return err.error
}

func (err *value) ErrorVerbose() string {
	return "value " + err.key + " = " + VerboseStringer(err.val)
}

func (err *value) Value() (key string, val any) {
	return err.key, err.val
}

// Get returns the values added to an error.
func Get(err error) map[string]any {
	vals := make(map[string]any)
	erriter.Iter(err, func(err error) {
		errv, ok := err.(interface { //nolint:errorlint // We want to compare the current error.
			Value() (key string, val any)
		})
		if !ok {
			return
		}
		k, v := errv.Value()
		_, ok = vals[k]
		if ok {
			return
		}
		vals[k] = v
	})
	return vals
}
