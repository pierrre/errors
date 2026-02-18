// Package errors provides error management.
//
// By convention, wrapping functions return a nil error if the provided error is nil.
package errors

import (
	std_errors "errors"
	"runtime"
	"strings"
	"testing"

	"github.com/pierrre/errors/errbase"
	"github.com/pierrre/errors/errmsg"
	"github.com/pierrre/errors/errstack"
)

// ErrUnsupported is an alias for [std_errors.ErrUnsupported].
var ErrUnsupported = std_errors.ErrUnsupported

// New returns a new error with a message and a stack.
//
// Warning: don't use this function to create a global (sentinel) error, as it will contain the stack of the (main) goroutine creating it.
// Use [errbase.New] instead.
func New(msg string) error {
	err := errbase.New(msg)
	err = errstack.WrapSkip(err, 1)
	if ReportGlobalInit != nil {
		checkGlobalInit(err, ReportGlobalInit)
	}
	return err
}

// Newf returns a new error with a formatted message and a stack.
//
// It supports the %w verb.
//
// Warning: don't use this function to create a global (sentinel) error, as it will contain the stack of the (main) goroutine creating it.
// Use [errbase.Newf] instead.
func Newf(format string, args ...any) error {
	err := errbase.Newf(format, args...)
	err = errstack.WrapSkip(err, 1)
	if ReportGlobalInit != nil {
		checkGlobalInit(err, ReportGlobalInit)
	}
	return err
}

// ReportGlobalInit reports a global error initialization.
// It is discouraged to call [New] or [Newf] to create a global (sentinel) error, as it will contain the stack of the (main) goroutine that created it.
// Instead, call [errbase.New] or [errbase.Newf], and [Wrap] it before returning it, which will add the stack of the goroutine returning the error.
//
// Example:
//
//	var ErrGlobal = errbase.New("global error")
//
//	func myFunc() error {
//		return errors.Wrap(ErrGlobal, "myFunc error")
//	}
//
// The default values's behavior is to panic during tests, and do nothing during normal execution.
// It can be disabled by setting it to nil.
//
// The implementation of [New] and [Newf] checks if the error is created by a function named "init".
// It doesn't report errors created by "init()" functions, which are named "init.N" where N is a number.
var ReportGlobalInit func(error) = func() func(error) {
	var f func(error)
	if testing.Testing() {
		f = func(err error) {
			panic(err)
		}
	}
	return f
}()

func checkGlobalInit(err error, report func(error)) {
	// This code doesn't call [errstack.Frames] to avoid memory allocations.
	errf, ok := err.(interface {
		StackFrames() []uintptr
	})
	if !ok {
		return
	}
	pcs := errf.StackFrames()
	if len(pcs) == 0 {
		return
	}
	f := runtime.FuncForPC(pcs[0])
	if !strings.HasSuffix(f.Name(), ".init") {
		return
	}
	err = Wrap(err, "global error initialization detected, use errbase.New() instead, see https://pkg.go.dev/github.com/pierrre/errors#ReportGlobalInit ")
	report(err)
}

// Wrap adds a message to an error, and a stack if it doesn't have one.
func Wrap(err error, msg string) error {
	if err != nil {
		err = errstack.EnsureSkip(err, 1)
		err = errmsg.Wrap(err, msg)
	}
	return err
}

// Wrapf adds a formatted message to an error, and a stack if it doesn't have one.
//
// It doesn't support the %w verb.
func Wrapf(err error, format string, args ...any) error {
	if err != nil {
		err = errstack.EnsureSkip(err, 1)
		err = errmsg.Wrapf(err, format, args...)
	}
	return err
}

// As is an alias for [std_errors.As].
func As(err error, target any) bool {
	return std_errors.As(err, target)
}

// AsType is an alias for [std_errors.AsType].
func AsType[E error](err error) (E, bool) {
	return std_errors.AsType[E](err)
}

// Is is an alias for [std_errors.Is].
func Is(err, target error) bool {
	return std_errors.Is(err, target)
}

// Join calls [std_errors.Join] and adds a stack.
func Join(errs ...error) error {
	err := std_errors.Join(errs...)
	err = errstack.WrapSkip(err, 1)
	return err
}

// Unwrap is an alias for [std_errors.Unwrap].
func Unwrap(err error) error {
	return std_errors.Unwrap(err)
}
