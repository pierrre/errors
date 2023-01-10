package errors

import (
	"fmt"
	"io"
	"strings"
)

// Verboser is an error that provides verbose information.
//
// It is used by Verbose().
type Verboser interface {
	// ErrorVerbose writes the verbose message of the error to the writer.
	// It must only write the verbose message of the error, not the error chain.
	// It is responsible for writing a new line character at the end.
	ErrorVerbose(io.Writer)
}

// Verbose writes the error's verbose message to the writer.
//
// The first line is the error's message.
// The following lines are the verbose message of the error chain.
func Verbose(w io.Writer, err error) {
	verbose(w, err, nil)
}

func verbose(w io.Writer, err error, depth []int) {
	verboseSub(w, depth)
	_, _ = fmt.Fprint(w, err)
	_, _ = io.WriteString(w, "\n")
	for ; err != nil; err = verboseNext(w, err, depth) {
		v, ok := err.(Verboser) //nolint:errorlint // We want to compare the current error.
		if ok {
			v.ErrorVerbose(w)
		}
	}
}

func verboseSub(w io.Writer, depth []int) {
	if len(depth) == 0 {
		return
	}
	_, _ = io.WriteString(w, "\nSub error ")
	for i, d := range depth {
		if i > 0 {
			_, _ = io.WriteString(w, ".")
		}
		_, _ = fmt.Fprint(w, d)
	}
	_, _ = io.WriteString(w, ": ")
}

func verboseNext(w io.Writer, err error, depth []int) error {
	switch err := err.(type) { //nolint: errorlint // We want to compare the current error.
	case interface{ Unwrap() error }:
		return err.Unwrap() //nolint:wrapcheck // We want to return the wrapped error.
	case interface{ Unwrap() []error }:
		errs := err.Unwrap()
		for i, err := range errs {
			verbose(w, err, append(depth, i))
		}
	}
	return nil
}

// VerboseString returns the error's verbose message as a string.
func VerboseString(err error) string {
	var b strings.Builder // TODO use a buffer pool.
	Verbose(&b, err)
	return b.String()
}

// VerboseFormatter returns a fmt.Formatter that writes the error's verbose message.
func VerboseFormatter(err error) fmt.Formatter {
	return &verboseFormatter{
		error: err,
	}
}

type verboseFormatter struct {
	error error
}

func (f *verboseFormatter) Format(s fmt.State, verb rune) {
	Verbose(s, f.error)
}
