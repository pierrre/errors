// Package errtag provides a way to add tags to errors.
package errtag

import (
	"iter"
	"strconv"

	"github.com/pierrre/errors/errappend"
	"github.com/pierrre/errors/erriter"
)

// Wrap adds a tag to an error.
//
// Tags should be used for short and simple values, such as identifiers.
//
// The verbose message is "tag <key> = <val>".
func Wrap(err error, key string, val string) error {
	if err == nil {
		return nil
	}
	return &tag{
		error: err,
		key:   key,
		val:   val,
	}
}

// WrapInt is a helper for [Wrap] with int value.
func WrapInt(err error, key string, value int) error {
	return Wrap(err, key, strconv.Itoa(value))
}

// WrapInt64 is a helper for [Wrap] with int64 value.
func WrapInt64(err error, key string, value int64) error {
	return Wrap(err, key, strconv.FormatInt(value, 10))
}

// WrapFloat64 is a helper for [Wrap] with float64 value.
func WrapFloat64(err error, key string, value float64) error {
	return Wrap(err, key, strconv.FormatFloat(value, 'g', -1, 64))
}

// WrapBool is a helper for [Wrap] with bool value.
func WrapBool(err error, key string, value bool) error {
	return Wrap(err, key, strconv.FormatBool(value))
}

type tag struct {
	error
	key string
	val string
}

func (err *tag) Unwrap() error {
	return err.error
}

func (err *tag) ErrorAppend(b []byte) []byte {
	return errappend.Append(b, err.error)
}

func (err *tag) ErrorVerboseAppend(b []byte) []byte {
	b = append(b, "tag "...)
	b = append(b, err.key...)
	b = append(b, " = "...)
	b = append(b, err.val...)
	return b
}

func (err *tag) Tag() (key string, val string) {
	return err.key, err.val
}

// All returns a [iter.Seq2] of tags added to an error.
func All(err error) iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		for err := range erriter.All(err) {
			errt, ok := err.(interface {
				Tag() (key string, val string)
			})
			if !ok {
				continue
			}
			k, v := errt.Tag()
			if !yield(k, v) {
				return
			}
		}
	}
}

// Get returns the tags added to an error.
func Get(err error) map[string]string {
	return erriter.FirstKeys(All(err))
}
