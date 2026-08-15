// Package errappend provides an interface for error types that want to provide a custom error message when appended to a byte slice.
package errappend

import (
	"github.com/pierrre/go-libs/bytesutil"
)

// Interface is an interface that can be implemented by error types that want to provide a custom error message when appended to a byte slice.
type Interface interface {
	error
	ErrorAppend(b []byte) []byte
}

// String returns the error message of the given error as a string.
func String(err Interface) string {
	if err == nil {
		return ""
	}
	bw := bytesWriterBool.Get()
	defer bytesWriterBool.Put(bw)
	*bw = err.ErrorAppend(*bw)
	return bw.String()
}

// Append appends the error message of the given error to the given byte slice and returns the resulting byte slice.
// If the given error is nil, the given byte slice is returned unchanged.
// If the given error implements the [Interface] interface, its [Interface.ErrorAppend] method is called to append the error message to the byte slice.
// Otherwise, the error message is appended to the byte slice using the [error.Error] method.
func Append(b []byte, err error) []byte {
	if err == nil {
		return b
	}
	if a, ok := err.(Interface); ok {
		return a.ErrorAppend(b)
	}
	return append(b, err.Error()...)
}

var bytesWriterBool = &bytesutil.WriterPool{}
