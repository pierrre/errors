package errmsg_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/errors"
	"github.com/pierrre/errors/errbase"
	. "github.com/pierrre/errors/errmsg"
)

func Example() {
	err := errbase.New("error")
	err = Wrap(err, "message")
	fmt.Println(err)
	// Output: message: error
}

func Test(t *testing.T) {
	err := errbase.New("error")
	err = Wrapf(err, "test %d", 1)
	assert.ErrorEqual(t, err, "test 1: error")
}

func TestNil(t *testing.T) {
	err := Wrap(nil, "test")
	assert.NoError(t, err)
}

func TestEmpty(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err, "")
	assert.ErrorEqual(t, err, "error")
}

func TestUnwrap(t *testing.T) {
	err1 := errbase.New("error")
	err2 := Wrap(err1, "test")
	err2 = errors.Unwrap(err2)
	assert.Equal(t, err2, err1)
}

func TestWrapAllocs(t *testing.T) {
	err := errbase.New("error")
	var res error
	assert.AllocsPerRun(t, 100, func() {
		res = Wrap(err, "test")
	}, 2)
	runtime.KeepAlive(res)
}

func TestWrapfAllocs(t *testing.T) {
	err := errbase.New("error")
	var res error
	assert.AllocsPerRun(t, 100, func() {
		res = Wrapf(err, "test %d", 1)
	}, 3)
	runtime.KeepAlive(res)
}

func TestErrorAllocs(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err, "test")
	var res string
	assert.AllocsPerRun(t, 100, func() {
		res = err.Error()
	}, 0)
	runtime.KeepAlive(res)
}

func BenchmarkWrap(b *testing.B) {
	err := errbase.New("error")
	var res error
	for b.Loop() {
		res = Wrap(err, "test")
	}
	runtime.KeepAlive(res)
}

func BenchmarkWrapf(b *testing.B) {
	err := errbase.New("error")
	var res error
	for b.Loop() {
		res = Wrapf(err, "test %d", 1)
	}
	runtime.KeepAlive(res)
}

func BenchmarkError(b *testing.B) {
	err := errbase.New("error")
	err = Wrap(err, "test")
	var res string
	for b.Loop() {
		res = err.Error()
	}
	runtime.KeepAlive(res)
}
