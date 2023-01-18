// Package errverbose provides utilities to manage error verbose messages.
package errverbose

import (
	"fmt"
	"io"
	"strings"

	"github.com/pierrre/errors"
)

// Interface is an error that provides verbose information.
//
// It is used by Write().
type Interface interface {
	// ErrorVerbose returns the error verbose message.
	// It must only return the verbose message of the error, not the error chain.
	ErrorVerbose() string
}

// Write writes the error's verbose message to the writer.
//
// The first line is the error's message.
// The following lines are the verbose message of the error chain.
func Write(w io.Writer, err error) {
	_, _ = fmt.Fprint(w, err)
	_, _ = io.WriteString(w, "\n")
	for ; err != nil; err = errors.Unwrap(err) {
		v, ok := err.(Interface) //nolint:errorlint // We want to compare the current error.
		if ok {
			s := v.ErrorVerbose()
			_, _ = io.WriteString(w, s)
			_, _ = io.WriteString(w, "\n")
		}
	}
}

// String returns the error's verbose message as a string.
func String(err error) string {
	var b strings.Builder // TODO use a buffer pool.
	Write(&b, err)
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
