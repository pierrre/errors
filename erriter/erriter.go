// Package erriter provides a way to iterate over an error tree.
package erriter

import (
	"iter"
	"slices"

	"github.com/pierrre/go-libs/syncutil"
)

// op is one step of the traversal performed by iterFunc.
type op struct {
	err  error
	exit bool // remove err from seen instead of visiting it
}

// walker holds the state of a single All traversal.
type walker struct {
	seen  map[error]struct{} // errors on the current root-to-node path
	stack []op
}

var walkerPool = &syncutil.Pool[*walker]{
	New: func() *walker {
		return &walker{
			seen: make(map[error]struct{}, 8),
		}
	},
}

func releaseWalker(w *walker) {
	clear(w.seen)
	clear(w.stack)
	w.stack = w.stack[:0]
	walkerPool.Put(w)
}

// All returns an [iter.Seq] that iterates over an error tree.
func All(err error) iter.Seq[error] {
	return func(yield func(error) bool) {
		iterFunc(err, yield)
	}
}

// iterFunc iterates over the error tree rooted at root in depth-first order: an error, then its multi errors in order, then its chain error.
// It stops as soon as f returns false.
// An error that is its own ancestor is not visited again, and iteration continues with the next sibling.
func iterFunc(root error, f func(err error) bool) {
	if root == nil {
		return
	}
	w := walkerPool.Get()
	defer releaseWalker(w)
	w.stack = append(w.stack, op{err: root})
	for len(w.stack) > 0 {
		last := len(w.stack) - 1
		o := w.stack[last]
		w.stack[last] = op{}
		w.stack = w.stack[:last]
		if o.exit {
			delete(w.seen, o.err)
			continue
		}
		if o.err == nil {
			continue
		}
		if _, ok := w.seen[o.err]; ok {
			continue // cycle
		}
		w.seen[o.err] = struct{}{}
		if !f(o.err) {
			return
		}
		errs, next := Unwrap(o.err)
		w.stack = append(w.stack, op{err: o.err, exit: true})
		if next != nil {
			w.stack = append(w.stack, op{err: next})
		}
		for _, e := range slices.Backward(errs) {
			w.stack = append(w.stack, op{err: e})
		}
	}
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
