// Package errverbose provides utilities to manage error verbose messages.
package errverbose

import (
	"fmt"
	"io"

	"github.com/pierrre/errors/erriter"
	"github.com/pierrre/go-libs/bufpool"
	"github.com/pierrre/go-libs/strconvio"
	"github.com/pierrre/go-libs/syncutil"
	"github.com/pierrre/go-libs/unsafeio"
)

// Interface is an error that provides verbose information.
//
// It is used by [Write].
type Interface interface {
	// ErrorVerbose writes the error verbose message.
	// It must only write the verbose message of the error, not the error chain.
	ErrorVerbose(w io.Writer)
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
	write(w, err, depth)
}

func write(w io.Writer, err error, depth []int) {
	writeSub(w, depth)
	if err == nil {
		_, _ = unsafeio.WriteString(w, "<nil>\n")
		return
	}
	_, _ = unsafeio.WriteString(w, err.Error())
	_, _ = unsafeio.WriteString(w, "\n")
	for ; err != nil; err = writeNext(w, err, depth) {
		v, ok := err.(Interface)
		if ok {
			v.ErrorVerbose(w)
			_, _ = unsafeio.WriteString(w, "\n")
		}
	}
}

func writeSub(w io.Writer, depth []int) {
	if len(depth) == 0 {
		return
	}
	_, _ = unsafeio.WriteString(w, "\nSub error ")
	for i, d := range depth {
		if i > 0 {
			_, _ = unsafeio.WriteString(w, ".")
		}
		_, _ = strconvio.WriteInt(w, int64(d), 10)
	}
	_, _ = unsafeio.WriteString(w, ": ")
}

func writeNext(w io.Writer, err error, depth []int) error {
	errs, err := erriter.Unwrap(err)
	for i, e := range errs {
		write(w, e, append(depth, i))
	}
	return err
}

var bufferPool = bufpool.Pool{}

// String returns the error's verbose message as a string.
func String(err error) string {
	b := bufferPool.Get()
	defer bufferPool.Put(b)
	Write(b, err)
	return b.String()
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
