package errstack_test

import (
	"fmt"
	"io"
	"runtime"
	"strings"
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/errors"
	"github.com/pierrre/errors/errbase"
	. "github.com/pierrre/errors/errstack"
	"github.com/pierrre/errors/errverbose"
)

func Example() {
	err := errors.New("error")
	err = Wrap(err)
	fmt.Println(err)
	sfs := Frames(err)
	fmt.Println(len(sfs))
	// Output:
	// error
	// 2
}

func Test(t *testing.T) {
	err := errbase.New("error")
	err = Ensure(err)
	err = Ensure(err)
	sfs := Frames(err)
	assert.SliceLen(t, sfs, 1)
	sf := sfs[0]
	assert.NotZero(t, sf)
	f, _ := sf.Next()
	assert.Equal(t, f.Function, "github.com/pierrre/errors/errstack_test.Test")
}

func TestNil(t *testing.T) {
	err := Wrap(nil)
	assert.NoError(t, err)
}

func TestError(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err)
	assert.ErrorEqual(t, err, "error")
}

func TestVerbose(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err)
	var v errverbose.Interface
	assert.ErrorAs(t, err, &v)
	sb := new(strings.Builder)
	v.ErrorVerbose(sb)
	s := sb.String()
	assert.RegexpMatch(t, `^stack\n(\t.+ .+:\d+\n)+$`, s)
}

func TestStackFrames(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err)
	var sErr interface {
		StackFrames() []uintptr
	}
	assert.ErrorAs(t, err, &sErr)
	pcs := sErr.StackFrames()
	assert.SliceNotEmpty(t, pcs)
}

func TestJoin(t *testing.T) {
	err := Wrap(
		errors.Join(
			Wrap(
				errbase.New("error 1"),
			),
			Wrap(
				errbase.New("error 2"),
			),
		),
	)
	sfs := Frames(err)
	assert.SliceLen(t, sfs, 3)
}

func TestWrapAllocs(t *testing.T) {
	err := errbase.New("error")
	var res error
	assert.AllocsPerRun(t, 100, func() {
		res = Wrap(err)
	}, 2)
	runtime.KeepAlive(res)
}

func TestEnsureAllocs(t *testing.T) {
	err := errbase.New("error")
	err = Ensure(err)
	var res error
	assert.AllocsPerRun(t, 100, func() {
		res = Ensure(err)
	}, 0)
	runtime.KeepAlive(res)
}

func TestFramesAllocs(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err)
	var res []*runtime.Frames
	assert.AllocsPerRun(t, 100, func() {
		res = Frames(err)
	}, 2)
	runtime.KeepAlive(res)
}

func TestVerboseAllocs(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err)
	var v errverbose.Interface
	assert.ErrorAs(t, err, &v)
	assert.AllocsPerRun(t, 100, func() {
		v.ErrorVerbose(io.Discard)
	}, 1)
}

func BenchmarkWrap(b *testing.B) {
	err := errbase.New("error")
	var res error
	for i := 0; i < b.N; i++ {
		res = Wrap(err)
	}
	runtime.KeepAlive(res)
}

func BenchmarkEnsure(b *testing.B) {
	err := errbase.New("error")
	err = Ensure(err)
	var res error
	for i := 0; i < b.N; i++ {
		res = Ensure(err)
	}
	runtime.KeepAlive(res)
}

func BenchmarkFrames(b *testing.B) {
	err := errbase.New("error")
	err = Wrap(err)
	var res []*runtime.Frames
	for i := 0; i < b.N; i++ {
		res = Frames(err)
	}
	runtime.KeepAlive(res)
}

func BenchmarkVerbose(b *testing.B) {
	err := errbase.New("error")
	err = Wrap(err)
	var v errverbose.Interface
	assert.ErrorAs(b, err, &v)
	for i := 0; i < b.N; i++ {
		v.ErrorVerbose(io.Discard)
	}
}
