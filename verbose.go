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
	ErrorVerbose() string
}

// Verbose writes the error's verbose message to the writer.
//
// The first line is the error's message.
// The following lines are the verbose message of the error chain.
func Verbose(w io.Writer, err error) {
	_, _ = fmt.Fprint(w, err)
	_, _ = io.WriteString(w, "\n")
	for ; err != nil; err = Unwrap(err) {
		v, ok := err.(Verboser) //nolint:errorlint // We want to compare the current error.
		if ok {
			s := v.ErrorVerbose()
			_, _ = io.WriteString(w, s)
			_, _ = io.WriteString(w, "\n")
		}
	}
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
