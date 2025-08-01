// Package errstack provides utilities to manage error stack traces.
package errstack

import (
	std_errors "errors" // Prevent import cycle.
	"io"
	"iter"
	"runtime"

	"github.com/pierrre/errors/errbase"
	"github.com/pierrre/errors/erriter"
	"github.com/pierrre/go-libs/runtimeutil"
	"github.com/pierrre/go-libs/unsafeio"
)

// Wrap adds a stack to an error.
//
// The verbose message contains the stack.
//
// See [runtime.Frames].
func Wrap(err error) error {
	return WrapSkip(err, 1)
}

// WrapSkip calls [Wrap], skipping the given number of frames.
func WrapSkip(err error, skip int) error {
	if err == nil {
		return nil
	}
	return &stack{
		error:   err,
		callers: runtimeutil.GetCallers(skip + 1),
	}
}

// Ensure adds a stack to an error if it does not already have one.
func Ensure(err error) error {
	return EnsureSkip(err, 1)
}

// EnsureSkip adds a stack to an error if it does not already have one, skipping the given number of frames.
func EnsureSkip(err error, skip int) error {
	if err != nil && !has(err) {
		err = WrapSkip(err, skip+1)
	}
	return err
}

type stack struct {
	error
	callers []uintptr
}

func (err *stack) Unwrap() error {
	return err.error
}

func (err *stack) Is(target error) bool {
	return target == err || target == errHas
}

func (err *stack) ErrorVerbose(w io.Writer) {
	_, _ = unsafeio.WriteString(w, "stack:\n")
	_, _ = runtimeutil.WriteFrames(w, runtimeutil.GetCallersFrames(err.callers))
}

// StackFrames returns the list of PCs associated to the error.
//
// It exists and is named StackFrames in order to be compatible with the Sentry library, which expects this name.
//
// There is no stability guarantee for this method.
func (err *stack) StackFrames() []uintptr {
	return err.callers
}

// Frames returns the list of [runtime.Frame] associated to an error.
func Frames(err error) iter.Seq[iter.Seq[runtime.Frame]] {
	return func(yield func(iter.Seq[runtime.Frame]) bool) {
		for err := range erriter.All(err) {
			errf, ok := err.(interface {
				StackFrames() []uintptr
			})
			if ok {
				fs := runtimeutil.GetCallersFrames(errf.StackFrames())
				if !yield(fs) {
					return
				}
			}
		}
	}
}

var errHas = errbase.New("stack")

func has(err error) bool {
	return std_errors.Is(err, errHas)
}
