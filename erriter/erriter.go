// Package erriter provides a way to iterate over an error tree.
package erriter

import (
	"iter"

	"github.com/pierrre/go-libs/syncutil"
)

var seenPool = &syncutil.Pool[map[error]struct{}]{
	New: func() map[error]struct{} {
		return make(map[error]struct{}, 8)
	},
}

func releaseSeen(seen map[error]struct{}) {
	clear(seen)
	seenPool.Put(seen)
}

// All returns an [iter.Seq] that iterates over an error tree recursively.
func All(err error) iter.Seq[error] {
	return func(yield func(error) bool) {
		seen := seenPool.Get()
		defer releaseSeen(seen)
		iterFunc(err, yield, seen)
	}
}

func iterFunc(err error, f func(err error) bool, seen map[error]struct{}) bool {
	if err == nil {
		return true
	}
	if _, ok := seen[err]; ok {
		return true // cycle
	}
	seen[err] = struct{}{}
	defer delete(seen, err)

	if !f(err) {
		return false
	}
	errs, next := Unwrap(err)
	for _, e := range errs {
		if !iterFunc(e, f, seen) {
			return false
		}
	}
	return iterFunc(next, f, seen)
}

// Unwrap unwraps an error.
//
// If the error implements `Unwrap() error`, it returns the unwrapped error.
// If the error implements `Unwrap() []error`, it returns the unwrapped errors.
// Otherwise, it returns nil.
func Unwrap(err error) ([]error, error) {
	switch err := err.(type) { //nolint:errorlint // We want to check which interface is implemented by the current error.
	case interface{ Unwrap() error }:
		return nil, err.Unwrap()
	case interface{ Unwrap() []error }:
		return err.Unwrap(), nil
	}
	return nil, nil
}

// FirstKeys returns a map built from seq, where the first value is kept for each key.
// It may return a nil map if there is no value.
func FirstKeys[V any](seq iter.Seq2[string, V]) map[string]V {
	var m map[string]V
	for k, v := range seq {
		if _, ok := m[k]; ok {
			continue
		}
		if m == nil {
			m = make(map[string]V)
		}
		m[k] = v
	}
	return m
}
