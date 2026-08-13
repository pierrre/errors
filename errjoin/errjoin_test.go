package errjoin_test

import (
	"errors"
	"fmt"
	"io/fs"
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/errors/errbase"
	. "github.com/pierrre/errors/errjoin"
)

var testSink any

func Example() {
	err := Join(errbase.New("error 1"), errbase.New("error 2"))
	fmt.Println(err)
	// Output: error 1
	// error 2
}

func Test(t *testing.T) {
	err1 := errbase.New("error 1")
	err2 := errbase.New("error 2")
	err := Join(err1, err2)
	errUnwrap, ok := assert.Type[interface {
		Unwrap() []error
	}](t, err)
	assert.True(t, ok)
	errs := errUnwrap.Unwrap()
	assert.SliceEqual(t, errs, []error{err1, err2})
}

func TestNil(t *testing.T) {
	err := Join(nil, nil)
	assert.Zero(t, err)
}

func TestEmpty(t *testing.T) {
	err := Join()
	assert.Zero(t, err)
}

func TestSingle(t *testing.T) {
	err1 := errbase.New("error")
	err := Join(err1)
	assert.ErrorEqual(t, err, "error")
	errUnwrap, ok := assert.Type[interface {
		Unwrap() []error
	}](t, err)
	assert.True(t, ok)
	errs := errUnwrap.Unwrap()
	assert.SliceEqual(t, errs, []error{err1})
}

func TestError(t *testing.T) {
	err := Join(errbase.New("error 1"), errbase.New("error 2"))
	assert.ErrorEqual(t, err, "error 1\nerror 2")
}

func TestIs(t *testing.T) {
	err1 := errbase.New("error 1")
	err := Join(err1, errbase.New("error 2"))
	ok := errors.Is(err, err1)
	assert.True(t, ok)
}

func TestAs(t *testing.T) {
	err1 := &fs.PathError{Op: "op", Path: "path", Err: errbase.New("error 1")}
	err := Join(err1, errbase.New("error 2"))
	var pathError *fs.PathError
	ok := errors.As(err, &pathError)
	assert.True(t, ok)
}

func TestJoinAllocs(t *testing.T) {
	err1 := errbase.New("error 1")
	err2 := errbase.New("error 2")
	var res error
	assert.AllocsPerRun(t, 100, func() {
		res = Join(err1, err2)
	}, 2)
	testSink = res
}

func TestErrorAllocs(t *testing.T) {
	err := Join(errbase.New("error 1"), errbase.New("error 2"))
	var res string
	assert.AllocsPerRun(t, 100, func() {
		res = err.Error()
	}, 2)
	testSink = res
}

func TestErrorSingleAllocs(t *testing.T) {
	err := Join(errbase.New("error"))
	var res string
	assert.AllocsPerRun(t, 100, func() {
		res = err.Error()
	}, 0)
	testSink = res
}

func BenchmarkJoin(b *testing.B) {
	err1 := errbase.New("error 1")
	err2 := errbase.New("error 2")
	for b.Loop() {
		_ = Join(err1, err2)
	}
}

func BenchmarkError(b *testing.B) {
	err := Join(errbase.New("error 1"), errbase.New("error 2"))
	for b.Loop() {
		_ = err.Error()
	}
}

func BenchmarkErrorSingle(b *testing.B) {
	err := Join(errbase.New("error"))
	for b.Loop() {
		_ = err.Error()
	}
}
