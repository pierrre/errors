// Package errval provides a way to add values to errors.
package errval

import (
	"fmt"
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
	get(err, vals)
	return vals
}

func get(err error, vals map[string]any) {
	for ; err != nil; err = getNext(err, vals) {
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
}

func getNext(err error, vals map[string]any) error {
	switch err := err.(type) { //nolint:errorlint // We want to compare the current error.
	case interface{ Unwrap() error }:
		return err.Unwrap() //nolint:wrapcheck // We want to return the wrapped error.
	case interface{ Unwrap() []error }:
		for _, err := range err.Unwrap() {
			get(err, vals)
		}
	}
	return nil
}
