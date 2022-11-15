package errors

import (
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"sync"
)

// Stack adds a stack to an error.
//
// The verbose message contains the stack.
//
// See https://pkg.go.dev/runtime#Frames .
func Stack(err error) error {
	return stackSkip(err, 2)
}

func stackSkip(err error, skip int) error {
	if err == nil {
		return nil
	}
	return &stack{
		error:   err,
		callers: callers(skip + 1),
	}
}

type stack struct {
	error
	callers []uintptr
}

func (err *stack) Unwrap() error {
	return err.error
}

// StackFrameVerboseWriter writes a runtime.Frame to an error verbose message.
//
// It must write a new line character at the end.
//
// It can be changed in order to customize how runtime.Frame are formatted.
var StackFrameVerboseWriter = func(w io.Writer, f runtime.Frame) {
	_, file := filepath.Split(f.File)
	_, _ = fmt.Fprintf(w, "\t%s %s:%d\n", f.Function, file, f.Line)
}

func (err *stack) ErrorVerbose(w io.Writer) {
	_, _ = io.WriteString(w, "stack\n")
	fs := err.RuntimeStackFrames()
	for more := true; more; {
		var f runtime.Frame
		f, more = fs.Next()
		StackFrameVerboseWriter(w, f)
	}
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

// StackFrames returns the list of runtime.Frames associated to an error.
//
// See https://pkg.go.dev/runtime#Frames .
func StackFrames(err error) []*runtime.Frames {
	var fss []*runtime.Frames
	for ; err != nil; err = Unwrap(err) {
		err, ok := err.(*stack) //nolint:errorlint // We want to compare the current error.
		if ok {
			fs := err.RuntimeStackFrames()
			fss = append(fss, fs)
		}
	}
	return fss
}

func ensureStack(err error, skip int) error {
	if !hasStack(err) {
		err = stackSkip(err, skip+1)
	}
	return err
}

func hasStack(err error) bool {
	var werr *stack
	return As(err, &werr)
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
	n := runtime.Callers(skip+1, pc)
	pcRes := make([]uintptr, n)
	copy(pcRes, pc)
	return pcRes
}
