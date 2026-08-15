package errignore_test

import (
	"fmt"
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/errors"
	"github.com/pierrre/errors/errappend"
	"github.com/pierrre/errors/errbase"
	. "github.com/pierrre/errors/errignore"
	"github.com/pierrre/errors/errmsg"
	"github.com/pierrre/errors/errverbose"
)

var testSink any

func Example() {
	err := errbase.New("error")
	err = Wrap(err)
	ignored := Is(err)
	fmt.Println(ignored)
	// Output: true
}

func Test(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err)
	ignored := Is(err)
	assert.True(t, ignored)
}

func TestNil(t *testing.T) {
	err := Wrap(nil)
	assert.NoError(t, err)
}

func TestFalse(t *testing.T) {
	err := errbase.New("error")
	ignored := Is(err)
	assert.False(t, ignored)
}

func TestError(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err)
	assert.ErrorEqual(t, err, "error")
}

func TestVerbose(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err)
	v, _ := assert.ErrorAsType[errverbose.Interface](t, err)
	s := v.ErrorVerbose()
	assert.Equal(t, s, "ignored")
}

func TestErrorAppend(t *testing.T) {
	err := errmsg.Wrap(errbase.New("error"), "msg")
	err = Wrap(err)
	_, _ = assert.ErrorAsType[errappend.Interface](t, err)
	b := errappend.Append(nil, err)
	assert.Equal(t, string(b), "msg: error")
}

func TestUnwrap(t *testing.T) {
	err1 := errbase.New("error")
	err2 := Wrap(err1)
	err2 = errors.Unwrap(err2)
	assert.Equal(t, err2, err1)
}

func TestWrapAllocs(t *testing.T) {
	err := errbase.New("error")
	var res error
	assert.AllocsPerRun(t, 100, func() {
		res = Wrap(err)
	}, 1)
	testSink = res
}

func TestIsAllocs(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err)
	var res bool
	assert.AllocsPerRun(t, 100, func() {
		res = Is(err)
	}, 1)
	testSink = res
}

func TestVerboseAllocs(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err)
	v, _ := assert.ErrorAsType[errverbose.Interface](t, err)
	assert.AllocsPerRun(t, 100, func() {
		_ = v.ErrorVerbose()
	}, 0)
}

func BenchmarkWrap(b *testing.B) {
	err := errbase.New("error")
	for b.Loop() {
		_ = Wrap(err)
	}
}

func BenchmarkIs(b *testing.B) {
	err := errbase.New("error")
	err = Wrap(err)
	for b.Loop() {
		_ = Is(err)
	}
}

func BenchmarkVerbose(b *testing.B) {
	err := errbase.New("error")
	err = Wrap(err)
	v, _ := assert.ErrorAsType[errverbose.Interface](b, err)
	for b.Loop() {
		_ = v.ErrorVerbose()
	}
}
