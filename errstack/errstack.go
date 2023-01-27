// Package errstack provides utilities to manage error stack traces.
package errstack

import (
	std_errors "errors" // Prevent import cycle.
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// Wrap adds a stack to an error.
//
// The verbose message contains the stack.
//
// See https://pkg.go.dev/runtime#Frames .
func Wrap(err error) error {
	return WrapSkip(err, 1)
}

// WrapSkip adds a stack to an error, skipping the given number of frames.
//
// See Wrap().
func WrapSkip(err error, skip int) error {
	if err == nil {
		return nil
	}
	return &stack{
		error:   err,
		callers: callers(skip + 1),
	}
}

// Ensure adds a stack to an error if it does not already have one.
func Ensure(err error) error {
	return EnsureSkip(err, 1)
}

// EnsureSkip adds a stack to an error if it does not already have one, skipping the given number of frames.
func EnsureSkip(err error, skip int) error {
	if !has(err) {
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

func (err *stack) ErrorVerbose() string {
	b := new(strings.Builder) // TODO use a buffer pool.
	_, _ = b.WriteString("stack\n")
	fs := err.RuntimeStackFrames()
	for more := true; more; {
		var f runtime.Frame
		f, more = fs.Next()
		_, file := filepath.Split(f.File)
		_, _ = fmt.Fprintf(b, "\t%s %s:%d\n", f.Function, file, f.Line)
	}
	return b.String()
}

// StackFrames returns the list of PCs associated to the error.
//
// It exists and is named StackFrames in order to be compatible with the Sentry library, which expects this name.
//
// There is no stability guarantee for this method.
func (err *stack) StackFrames() []uintptr {
	return err.callers
}

// RuntimeStackFrames returns the runtime.Frames associated to the error.
//
// It should be named StackFrames, but it was not possible because of the compatibility with the Sentry library.
//
// There is no stability guarantee for this method.
func (err *stack) RuntimeStackFrames() *runtime.Frames {
	return runtime.CallersFrames(err.callers)
}

// Frames returns the list of runtime.Frames associated to an error.
//
// See https://pkg.go.dev/runtime#Frames .
func Frames(err error) []*runtime.Frames {
	var fss []*runtime.Frames
	for ; err != nil; err = std_errors.Unwrap(err) {
		err, ok := err.(interface { //nolint:errorlint // We want to compare the current error.
			RuntimeStackFrames() *runtime.Frames
		})
		if ok {
			fs := err.RuntimeStackFrames()
			fss = append(fss, fs)
		}
	}
	return fss
}

func has(err error) bool {
	var werr interface {
		RuntimeStackFrames() *runtime.Frames
	}
	return std_errors.As(err, &werr)
}

const callersMaxLength = 1 << 16

var callersPool = sync.Pool{
	New: func() any {
		return make([]uintptr, callersMaxLength)
	},
}

func callers(skip int) []uintptr {
	pcItf := callersPool.Get()
	defer callersPool.Put(pcItf)
	pc := pcItf.([]uintptr) //nolint:forcetypeassert // The pool always contains []uintptr.
	n := runtime.Callers(skip+2, pc)
	pcRes := make([]uintptr, n)
	copy(pcRes, pc)
	return pcRes
}
