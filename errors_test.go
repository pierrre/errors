package errors_test

import (
	"fmt"
	"io/fs"
	"runtime"
	"testing"

	"github.com/pierrre/assert"
	. "github.com/pierrre/errors"
	"github.com/pierrre/errors/errbase"
	"github.com/pierrre/errors/errstack"
)

func ExampleNew() {
	err := New("error")
	fmt.Println(err)
	// Output: error
}

func ExampleWrap() {
	err := errbase.New("error")
	err = Wrap(err, "wrap")
	fmt.Println(err)
	// Output: wrap: error
}

func TestNew(t *testing.T) {
	err := New("error")
	assert.ErrorEqual(t, err, "error")
	sfs := errstack.Frames(err)
	assert.SliceLen(t, sfs, 1)
}

func TestNewf(t *testing.T) {
	err := Newf("error %d", 1)
	assert.ErrorEqual(t, err, "error 1")
	sfs := errstack.Frames(err)
	assert.SliceLen(t, sfs, 1)
}

func TestReportGlobalInitPanics(t *testing.T) {
	assert.Panics(t, func() {
		ReportGlobalInit(New("error"))
	})
}

var errGlobal = errstack.Wrap(errbase.New("global error"))

func TestCheckGlobalInit(t *testing.T) {
	called := false
	CheckGlobalInit(errGlobal, func(err error) {
		called = true
		assert.ErrorEqual(t, err, "global error initialization detected, use errbase.New() instead, see https://pkg.go.dev/github.com/pierrre/errors#ReportGlobalInit : global error")
	})
	assert.True(t, called)
}

func TestCheckGlobalInitNoStack(t *testing.T) {
	err := errbase.New("error")
	CheckGlobalInit(err, nil)
}

func TestCheckGlobalInitEmptyPCs(t *testing.T) {
	err := errstack.WrapSkip(errbase.New("error"), 1000)
	CheckGlobalInit(err, nil)
}

func TestCheckGlobalInitNotInit(t *testing.T) {
	err := New("error")
	CheckGlobalInit(err, nil)
}

func TestWrap(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err, "test")
	assert.ErrorEqual(t, err, "test: error")
	sfs := errstack.Frames(err)
	assert.SliceLen(t, sfs, 1)
}

func TestWrapf(t *testing.T) {
	err := errbase.New("error")
	err = Wrapf(err, "test %d", 1)
	assert.ErrorEqual(t, err, "test 1: error")
	sfs := errstack.Frames(err)
	assert.SliceLen(t, sfs, 1)
}

func TestAs(t *testing.T) {
	err := errbase.New("error")
	err = &fs.PathError{Err: err}
	err = Wrap(err, "test")
	var pathError *fs.PathError
	ok := As(err, &pathError)
	assert.True(t, ok)
}

func TestIs(t *testing.T) {
	errBase := errbase.New("error")
	err := Wrap(errBase, "test")
	ok := Is(err, errBase)
	assert.True(t, ok)
}

func TestJoin(t *testing.T) {
	err := Join(New("error 1"), New("error 2"))
	err = Unwrap(err) // Remove the stack.
	errUnwrap, _ := assert.Type[interface {
		Unwrap() []error
	}](t, err)
	errs := errUnwrap.Unwrap()
	assert.SliceLen(t, errs, 2)
}

func TestUnwrap(t *testing.T) {
	errBase := errbase.New("error")
	err := Wrap(errBase, "test")
	err = Unwrap(err)
	err = Unwrap(err)
	assert.Equal(t, err, errBase)
}

func TestNewAllocs(t *testing.T) {
	var res error
	assert.AllocsPerRun(t, 100, func() {
		res = New("error")
	}, 3)
	runtime.KeepAlive(res)
}

func TestNewfAllocs(t *testing.T) {
	var res error
	assert.AllocsPerRun(t, 100, func() {
		res = Newf("error %d", 1)
	}, 4)
	runtime.KeepAlive(res)
}

func TestWrapAllocs(t *testing.T) {
	err := errbase.New("error")
	var res error
	assert.AllocsPerRun(t, 100, func() {
		res = Wrap(err, "test")
	}, 4)
	runtime.KeepAlive(res)
}

func TestWrapfAllocs(t *testing.T) {
	err := errbase.New("error")
	var res error
	assert.AllocsPerRun(t, 100, func() {
		res = Wrapf(err, "test %d", 1)
	}, 5)
	runtime.KeepAlive(res)
}

func BenchmarkNew(b *testing.B) {
	var res error
	for range b.N {
		res = New("error")
	}
	runtime.KeepAlive(res)
}

func BenchmarkNewf(b *testing.B) {
	var res error
	for range b.N {
		res = Newf("error %d", 1)
	}
	runtime.KeepAlive(res)
}

func BenchmarkWrap(b *testing.B) {
	err := errbase.New("error")
	var res error
	for range b.N {
		res = Wrap(err, "test")
	}
	runtime.KeepAlive(res)
}

func BenchmarkWrapf(b *testing.B) {
	err := errbase.New("error")
	var res error
	for range b.N {
		res = Wrapf(err, "test %d", 1)
	}
	runtime.KeepAlive(res)
}
