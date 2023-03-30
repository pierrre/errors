// Package errverbose provides utilities to manage error verbose messages.
package errverbose

import (
	"fmt"
	"io"

	"github.com/pierrre/go-libs/bufpool"
	"github.com/pierrre/go-libs/strconvio"
	"github.com/pierrre/go-libs/syncutil"
)

// Interface is an error that provides verbose information.
//
// It is used by Write().
type Interface interface {
	// ErrorVerbose returns the error verbose message.
	// It must only return the verbose message of the error, not the error chain.
	ErrorVerbose() string
}

var depthPool = syncutil.Pool[[]int]{
	New: func() *[]int {
		depth := make([]int, 100)
		return &depth
	},
}

// Write writes the error's verbose message to the writer.
//
// The first line is the error's message.
// The following lines are the verbose message of the error chain.
func Write(w io.Writer, err error) {
	depthP := depthPool.Get()
	defer depthPool.Put(depthP)
	depth := *depthP
	write(w, err, depth[:0])
}

func write(w io.Writer, err error, depth []int) {
	writeSub(w, depth)
	_, _ = io.WriteString(w, err.Error())
	_, _ = io.WriteString(w, "\n")
	for ; err != nil; err = writeNext(w, err, depth) {
		v, ok := err.(Interface) //nolint:errorlint // We want to compare the current error.
		if ok {
			s := v.ErrorVerbose()
			_, _ = io.WriteString(w, s)
			_, _ = io.WriteString(w, "\n")
		}
	}
}

func writeSub(w io.Writer, depth []int) {
	if len(depth) == 0 {
		return
	}
	_, _ = io.WriteString(w, "\nSub error ")
	for i, d := range depth {
		if i > 0 {
			_, _ = io.WriteString(w, ".")
		}
		_, _ = strconvio.WriteInt(w, int64(d), 10)
	}
	_, _ = io.WriteString(w, ": ")
}

func writeNext(w io.Writer, err error, depth []int) error {
	switch err := err.(type) { //nolint: errorlint // We want to compare the current error.
	case interface{ Unwrap() error }:
		return err.Unwrap() //nolint:wrapcheck // We want to return the wrapped error.
	case interface{ Unwrap() []error }:
		errs := err.Unwrap()
		for i, err := range errs {
			write(w, err, append(depth, i))
		}
	}
	return nil
}

var bufferPool = bufpool.Pool{}

// String returns the error's verbose message as a string.
func String(err error) string {
	b := bufferPool.Get()
	defer bufferPool.Put(b)
	b.Reset()
	Write(b, err)
	return b.String()
}

// Formatter returns a fmt.Formatter that writes the error's verbose message.
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
