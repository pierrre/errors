package errmsg_test

import (
	"fmt"
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/errors"
	"github.com/pierrre/errors/errbase"
	. "github.com/pierrre/errors/errmsg"
)

var testSink any

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
	testSink = res
}

func TestWrapfAllocs(t *testing.T) {
	err := errbase.New("error")
	var res error
	assert.AllocsPerRun(t, 100, func() {
		res = Wrapf(err, "test %d", 1)
	}, 3)
	testSink = res
}

func TestErrorAllocs(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err, "test")
	var res string
	assert.AllocsPerRun(t, 100, func() {
		res = err.Error()
	}, 0)
	testSink = res
}

func BenchmarkWrap(b *testing.B) {
	err := errbase.New("error")
	for b.Loop() {
		_ = Wrap(err, "test")
	}
}

func BenchmarkWrapf(b *testing.B) {
	err := errbase.New("error")
	for b.Loop() {
		_ = Wrapf(err, "test %d", 1)
	}
}

func BenchmarkError(b *testing.B) {
	err := errbase.New("error")
	err = Wrap(err, "test")
	for b.Loop() {
		_ = err.Error()
	}
}
