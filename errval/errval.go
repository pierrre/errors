// Package errval provides a way to add values to errors.
package errval

import (
	"io"
	"iter"

	"github.com/pierrre/errors/errappend"
	"github.com/pierrre/errors/erriter"
	"github.com/pierrre/go-libs/bytesutil"
	"github.com/pierrre/go-libs/syncutil/atomicutil"
	"github.com/pierrre/pretty"
)

// VerboseWriter writes the representation of a value in a verbose message.
//
// It can be changed in order to customize how values are formatted.
//
// By default it uses [pretty.Write].
var VerboseWriter atomicutil.Value[func(io.Writer, any)]

func init() {
	VerboseWriter.Store(prettyWrite)
}

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

func (err *value) ErrorAppend(b []byte) []byte {
	return errappend.Append(b, err.error)
}

func (err *value) ErrorVerboseAppend(b []byte) []byte {
	b = append(b, "value "...)
	b = append(b, err.key...)
	b = append(b, " = "...)
	bw := bytesWriterPool.Get()
	defer bytesWriterPool.Put(bw)
	VerboseWriter.Load()(bw, err.val)
	b = append(b, *bw...)
	return b
}

func (err *value) Value() (key string, val any) {
	return err.key, err.val
}

// All returns a [iter.Seq2] of values added to an error.
func All(err error) iter.Seq2[string, any] {
	return func(yield func(string, any) bool) {
		for err := range erriter.All(err) {
			errv, ok := err.(interface {
				Value() (key string, val any)
			})
			if !ok {
				continue
			}
			k, v := errv.Value()
			if !yield(k, v) {
				return
			}
		}
	}
}

// Get returns the values added to an error.
func Get(err error) map[string]any {
	return erriter.FirstKeys(All(err))
}

// GetValue returns the first value added to an error for the given key.
//
// It returns false if the key is not found.
func GetValue(err error, key string) (any, bool) {
	for k, v := range All(err) {
		if k == key {
			return v, true
		}
	}
	return nil, false
}

var bytesWriterPool = &bytesutil.WriterPool{}
