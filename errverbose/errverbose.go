// Package errverbose provides utilities to manage error verbose messages.
package errverbose

import (
	"fmt"
	"io"
	"strconv"

	"github.com/pierrre/errors/errappend"
	"github.com/pierrre/errors/erriter"
	"github.com/pierrre/go-libs/bytesutil"
	"github.com/pierrre/go-libs/syncutil"
)

// Interface is an error that provides verbose information.
//
// It is used by [Write].
//
// [AppendInterface] is an alternative that appends the verbose message to a byte slice instead of returning a string, avoiding the string allocation.
type Interface interface {
	error
	// ErrorVerbose returns the error verbose message.
	// It must only return the verbose message of the error, not the error chain.
	// It must not end with a newline.
	ErrorVerbose() string
}

// AppendInterface is an interface that can be implemented by error types that want to provide a custom verbose error message when appended to a byte slice.
//
// It is an alternative to [Interface] that appends the verbose message directly to the destination byte slice.
// This avoids the intermediate string allocation of [Interface], which is beneficial for performance.
type AppendInterface interface {
	error
	ErrorVerboseAppend(b []byte) []byte
}

var depthPool = syncutil.Pool[*[]int]{
	New: func() *[]int {
		v := make([]int, 100)
		return &v
	},
}

// Write writes the error's verbose message to the writer.
//
// The first line is the error's message.
// The following lines are the verbose message of the error chain.
func Write(w io.Writer, err error) {
	depthP := depthPool.Get()
	defer depthPool.Put(depthP)
	depth := (*depthP)[:0]
	bw, ok := w.(*bytesutil.Writer)
	if !ok {
		bw = bytesWriterPool.Get()
		defer func() {
			_, _ = w.Write(*bw)
			bytesWriterPool.Put(bw)
		}()
	}
	write(bw, err, depth)
}

func write(bw *bytesutil.Writer, err error, depth []int) {
	writeSub(bw, depth)
	if err == nil {
		bw.AppendString("<nil>\n")
		return
	}
	*bw = errappend.Append(*bw, err)
	bw.AppendByte('\n')
	for ; err != nil; err = writeNext(bw, err, depth) {
		switch v := err.(type) { //nolint:errorlint // We want to check for specific error types.
		case Interface:
			bw.AppendString(v.ErrorVerbose())
			bw.AppendByte('\n')
		case AppendInterface:
			*bw = v.ErrorVerboseAppend(*bw)
			bw.AppendByte('\n')
		}
	}
}

func writeSub(bw *bytesutil.Writer, depth []int) {
	if len(depth) == 0 {
		return
	}
	bw.AppendString("\nSub error ")
	for i, d := range depth {
		if i > 0 {
			bw.AppendString(".")
		}
		*bw = strconv.AppendInt(*bw, int64(d), 10)
	}
	bw.AppendString(": ")
}

func writeNext(bw *bytesutil.Writer, err error, depth []int) error {
	errs, err := erriter.Unwrap(err)
	for i, e := range errs {
		write(bw, e, append(depth, i))
	}
	return err
}

var bytesWriterPool = &bytesutil.WriterPool{}

// String returns the error's verbose message as a string.
func String(err error) string {
	bw := bytesWriterPool.Get()
	defer bytesWriterPool.Put(bw)
	Write(bw, err)
	return bw.String()
}

// Formatter returns a [fmt.Formatter] that writes the error's verbose message.
func Formatter(err error) fmt.Formatter {
	return &formatter{
		error: err,
	}
}

type formatter struct {
	error error
}

func (f *formatter) Format(s fmt.State, verb rune) {
	Write(s, f.error)
}
