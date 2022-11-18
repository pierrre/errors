// Package errval provides a way to add values to errors.
package errval

import (
	"fmt"
	"io"

	"github.com/pierrre/errors"
)

// VerboseWriter writes a value to an error verbose message.
//
// It must write a new line character at the end.
//
// It can be changed in order to customize how values are formatted.
var VerboseWriter = func(w io.Writer, v any) {
	_, _ = fmt.Fprint(w, v)
	_, _ = io.WriteString(w, "\n")
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

func (err *value) ErrorVerbose(w io.Writer) {
	_, _ = fmt.Fprintf(w, "value %s = ", err.key)
	VerboseWriter(w, err.val)
}

func (err *value) Value() (key string, val any) {
	return err.key, err.val
}

// Get returns the values added to an error.
func Get(err error) map[string]any {
	vals := make(map[string]any)
	for ; err != nil; err = errors.Unwrap(err) {
		err, ok := err.(interface { //nolint:errorlint // We want to compare the current error.
			Value() (key string, val any)
		})
		if !ok {
			continue
		}
		k, v := err.Value()
		_, ok = vals[k]
		if ok {
			continue
		}
		vals[k] = v
	}
	return vals
}
