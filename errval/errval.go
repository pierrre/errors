// Package errval provides a way to add values to errors.
package errval

import (
	"io"

	"github.com/pierrre/errors/erriter"
	"github.com/pierrre/go-libs/unsafeio"
	"github.com/pierrre/pretty"
)

// VerboseWriter writes the representation of a value in a verbose message.
//
// It can be changed in order to customize how values are formatted.
//
// By default it uses [pretty.Write].
var VerboseWriter func(io.Writer, any) = prettyWrite

func prettyWrite(w io.Writer, v any) {
	pretty.Write(w, v)
}

// Wrap adds a value to an error.
//
// The verbose message is "value <key> = <val>".
// The value is written using the [VerboseWriter] function.
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
	_, _ = unsafeio.WriteString(w, "value ")
	_, _ = unsafeio.WriteString(w, err.key)
	_, _ = unsafeio.WriteString(w, " = ")
	VerboseWriter(w, err.val)
}

func (err *value) Value() (key string, val any) {
	return err.key, err.val
}

// Get returns the values added to an error.
func Get(err error) map[string]any {
	vals := make(map[string]any)
	for err := range erriter.All(err) {
		errv, ok := err.(interface {
			Value() (key string, val any)
		})
		if !ok {
			continue
		}
		k, v := errv.Value()
		_, ok = vals[k]
		if ok {
			continue
		}
		vals[k] = v
	}
	return vals
}
