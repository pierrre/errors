// Package errstack provides utilities to manage error stack traces.
package errstack

import (
	std_errors "errors" // Prevent import cycle.
	"path/filepath"
	"runtime"

	"github.com/pierrre/errors/internal/strconvio"
	"github.com/pierrre/go-libs/bufpool"
	"github.com/pierrre/go-libs/syncutil"
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

var bufferPool = bufpool.Pool{}

func (err *stack) ErrorVerbose() string {
	b := bufferPool.Get()
	defer bufferPool.Put(b)
	b.Reset()
	_, _ = b.WriteString("stack\n")
	fs := err.RuntimeStackFrames()
	for more := true; more; {
		var f runtime.Frame
		f, more = fs.Next()
		_, file := filepath.Split(f.File)
		_, _ = b.WriteString("\t")
		_, _ = b.WriteString(f.Function)
		_, _ = b.WriteString(" ")
		_, _ = b.WriteString(file)
		_, _ = b.WriteString(":")
		_, _ = strconvio.WriteInt(b, int64(f.Line), 10)
		_, _ = b.WriteString("\n")
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
	for ; err != nil; err = stackFramesNext(err, &fss) {
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

func stackFramesNext(err error, pfss *[]*runtime.Frames) error {
	switch err := err.(type) { //nolint:errorlint // We want to compare the current error.
	case interface{ Unwrap() error }:
		return err.Unwrap() //nolint:wrapcheck // We want to return the wrapped error.
	case interface{ Unwrap() []error }:
		for _, err := range err.Unwrap() {
			fss := Frames(err)
			*pfss = append(*pfss, fss...)
		}
	}
	return nil
}

func has(err error) bool {
	var werr interface {
		RuntimeStackFrames() *runtime.Frames
	}
	return std_errors.As(err, &werr)
}

const callersMaxLength = 1 << 16

var callersPool = syncutil.Pool[[]uintptr]{
	New: func() *[]uintptr {
		pc := make([]uintptr, callersMaxLength)
		return &pc
	},
}

func callers(skip int) []uintptr {
	pcp := callersPool.Get()
	defer callersPool.Put(pcp)
	pc := *pcp
	n := runtime.Callers(skip+2, pc)
	pcRes := make([]uintptr, n)
	copy(pcRes, pc)
	return pcRes
}
