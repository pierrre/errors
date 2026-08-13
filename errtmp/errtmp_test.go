package errtmp_test

import (
	"fmt"
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/errors"
	"github.com/pierrre/errors/errappend"
	"github.com/pierrre/errors/errbase"
	"github.com/pierrre/errors/errmsg"
	. "github.com/pierrre/errors/errtmp"
	"github.com/pierrre/errors/errverbose"
)

var testSink any

func Example() {
	err := errbase.New("error")
	err = Wrap(err, true)
	temporary := Is(err)
	fmt.Println(temporary)
	// Output: true
}

func TestTrue(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err, true)
	temporary := Is(err)
	assert.True(t, temporary)
}

func TestFalse(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err, false)
	temporary := Is(err)
	assert.False(t, temporary)
}

func TestDefault(t *testing.T) {
	err := errbase.New("error")
	temporary := Is(err)
	assert.True(t, temporary)
}

func TestNil(t *testing.T) {
	err := Wrap(nil, true)
	assert.NoError(t, err)
}

func TestError(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err, true)
	assert.ErrorEqual(t, err, "error")
}

func TestVerbose(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err, true)
	var v errverbose.AppendInterface
	assert.ErrorAs(t, err, &v)
	b := v.ErrorVerboseAppend(nil)
	assert.Equal(t, string(b), "temporary = true")
}

func TestUnwrap(t *testing.T) {
	err1 := errbase.New("error")
	err2 := Wrap(err1, true)
	err2 = errors.Unwrap(err2)
	assert.Equal(t, err2, err1)
}

func TestErrorAppend(t *testing.T) {
	err := errmsg.Wrap(errbase.New("error"), "msg")
	err = Wrap(err, true)
	var i errappend.Interface
	assert.ErrorAs(t, err, &i)
	b := errappend.Append(nil, err)
	assert.Equal(t, string(b), "msg: error")
}

func TestWrapAllocs(t *testing.T) {
	err := errbase.New("error")
	var res error
	assert.AllocsPerRun(t, 100, func() {
		res = Wrap(err, true)
	}, 1)
	testSink = res
}

func TestIsAllocs(t *testing.T) {
	err := errbase.New("error")
	var res bool
	assert.AllocsPerRun(t, 100, func() {
		res = Is(err)
	}, 1)
	testSink = res
}

func TestVerboseAllocs(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err, true)
	var v errverbose.AppendInterface
	assert.ErrorAs(t, err, &v)
	var b []byte
	assert.AllocsPerRun(t, 100, func() {
		b = v.ErrorVerboseAppend(b)
		b = b[:0]
	}, 0)
}

func BenchmarkWrap(b *testing.B) {
	err := errbase.New("error")
	for b.Loop() {
		_ = Wrap(err, true)
	}
}

func BenchmarkIs(b *testing.B) {
	err := errbase.New("error")
	err = Wrap(err, true)
	for b.Loop() {
		_ = Is(err)
	}
}

func BenchmarkVerbose(b *testing.B) {
	err := errbase.New("error")
	err = Wrap(err, true)
	var v errverbose.AppendInterface
	assert.ErrorAs(b, err, &v)
	var buf []byte
	for b.Loop() {
		buf = v.ErrorVerboseAppend(buf)
		buf = buf[:0]
	}
}
