// Package errtag provides a way to add tags to errors.
package errtag

import (
	"strconv"
)

// Wrap adds a tag to an error.
//
// Tags should be use for short and simple values, such as identifiers.
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

// WrapInt is a helper for Wrap with int value.
func WrapInt(err error, key string, value int) error {
	return Wrap(err, key, strconv.Itoa(value))
}

// WrapInt64 is a helper for Wrap with int64 value.
func WrapInt64(err error, key string, value int64) error {
	return Wrap(err, key, strconv.FormatInt(value, 10))
}

// WrapFloat64 is a helper for Wrap with float64 value.
func WrapFloat64(err error, key string, value float64) error {
	return Wrap(err, key, strconv.FormatFloat(value, 'g', -1, 64))
}

// WrapBool is a helper for Wrap with bool value.
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

func (err *tag) ErrorVerbose() string {
	return "tag " + err.key + " = " + err.val
}

func (err *tag) Tag() (key string, val string) {
	return err.key, err.val
}

// Get returns the tags added to an error.
func Get(err error) map[string]string {
	tags := make(map[string]string)
	get(err, tags)
	return tags
}

func get(err error, tags map[string]string) {
	for ; err != nil; err = getNext(err, tags) {
		err, ok := err.(interface { //nolint:errorlint // We want to compare the current error.
			Tag() (key string, val string)
		})
		if !ok {
			continue
		}
		k, v := err.Tag()
		_, ok = tags[k]
		if ok {
			continue
		}
		tags[k] = v
	}
}

func getNext(err error, tags map[string]string) error {
	switch err := err.(type) { //nolint:errorlint // We want to compare the current error.
	case interface{ Unwrap() error }:
		return err.Unwrap() //nolint:wrapcheck // We want to return the wrapped error.
	case interface{ Unwrap() []error }:
		for _, err := range err.Unwrap() {
			get(err, tags)
		}
	}
	return nil
}
