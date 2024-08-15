// Package errstack provides utilities to manage error stack traces.
package errstack

import (
	std_errors "errors" // Prevent import cycle.
	"io"
	"path/filepath"
	"runtime"

	"github.com/pierrre/errors/errbase"
	"github.com/pierrre/errors/erriter"
	"github.com/pierrre/go-libs/strconvio"
	"github.com/pierrre/go-libs/syncutil"
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

func (err *stack) Is(target error) bool {
	return target == err || target == errHas
}

func (err *stack) ErrorVerbose(w io.Writer) {
	_, _ = unsafeio.WriteString(w, "stack\n")
	fs := err.RuntimeStackFrames()
	for more := true; more; {
		var f runtime.Frame
		f, more = fs.Next()
		_, file := filepath.Split(f.File)
		_, _ = unsafeio.WriteString(w, "\t")
		_, _ = unsafeio.WriteString(w, f.Function)
		_, _ = unsafeio.WriteString(w, " ")
		_, _ = unsafeio.WriteString(w, file)
		_, _ = unsafeio.WriteString(w, ":")
		_, _ = strconvio.WriteInt(w, int64(f.Line), 10)
		_, _ = unsafeio.WriteString(w, "\n")
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

// RuntimeStackFrames returns the [runtime.Frames] associated to the error.
//
// It should be named StackFrames, but it was not possible because of the compatibility with the Sentry library.
//
// There is no stability guarantee for this method.
func (err *stack) RuntimeStackFrames() *runtime.Frames {
	return runtime.CallersFrames(err.callers)
}

// Frames returns the list of [runtime.Frames] associated to an error.
func Frames(err error) []*runtime.Frames {
	var fss []*runtime.Frames
	erriter.Iter(err, func(err error) {
		errf, ok := err.(interface {
			RuntimeStackFrames() *runtime.Frames
		})
		if ok {
			fs := errf.RuntimeStackFrames()
			fss = append(fss, fs)
		}
	})
	return fss
}

var errHas = errbase.New("stack")

func has(err error) bool {
	return std_errors.Is(err, errHas)
}

const callersMaxLength = 1 << 16

var callersPool = syncutil.PoolFor[[]uintptr]{
	New: func() *[]uintptr {
		v := make([]uintptr, callersMaxLength)
		return &v
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
