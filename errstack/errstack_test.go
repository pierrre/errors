package errstack_test

import (
	"fmt"
	"iter"
	"runtime"
	"slices"
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/errors"
	"github.com/pierrre/errors/errappend"
	"github.com/pierrre/errors/errbase"
	"github.com/pierrre/errors/errmsg"
	. "github.com/pierrre/errors/errstack"
	"github.com/pierrre/errors/errverbose"
)

var testSink any

func Example() {
	err := errors.New("error")
	err = Wrap(err)
	fmt.Println(err)
	sfs := slices.Collect(Frames(err))
	fmt.Println(len(sfs))
	// Output:
	// error
	// 2
}

func Test(t *testing.T) {
	err := errbase.New("error")
	err = Ensure(err)
	err = Ensure(err)
	sfs := slices.Collect(Frames(err))
	assert.SliceLen(t, sfs, 1)
	sf := slices.Collect(sfs[0])
	assert.SliceNotEmpty(t, sf)
	f := sf[0]
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
	v, _ := assert.ErrorAsType[errverbose.AppendInterface](t, err)
	b := v.ErrorVerboseAppend(nil)
	s := string(b)
	t.Log(s)
	assert.RegexpMatch(t, `^stack:\n(.+\n\t.+:\d+\n)+$`, s)
}

func TestErrorAppend(t *testing.T) {
	err := errmsg.Wrap(errbase.New("error"), "msg")
	err = Wrap(err)
	_, _ = assert.ErrorAsType[errappend.Interface](t, err)
	b := errappend.Append(nil, err)
	assert.Equal(t, string(b), "msg: error")
}

func TestStackFrames(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err)
	sErr, _ := assert.ErrorAsType[interface {
		error
		StackFrames() []uintptr
	}](t, err)
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
	sfs := slices.Collect(Frames(err))
	assert.SliceLen(t, sfs, 4)
}

func TestWrapAllocs(t *testing.T) {
	err := errbase.New("error")
	var res error
	assert.AllocsPerRun(t, 100, func() {
		res = Wrap(err)
	}, 2)
	testSink = res
}

func TestEnsureAllocs(t *testing.T) {
	err := errbase.New("error")
	err = Ensure(err)
	var res error
	assert.AllocsPerRun(t, 100, func() {
		res = Ensure(err)
	}, 0)
	testSink = res
}

func TestFramesInterrupt(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err)
	for range Frames(err) {
		break
	}
}

func TestFramesAllocs(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err)
	var res iter.Seq[iter.Seq[runtime.Frame]]
	assert.AllocsPerRun(t, 100, func() {
		res = Frames(err)
	}, 1)
	testSink = res
}

func TestVerboseAllocs(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err)
	v, _ := assert.ErrorAsType[errverbose.AppendInterface](t, err)
	var b []byte
	assert.AllocsPerRun(t, 100, func() {
		b = v.ErrorVerboseAppend(b)
		b = b[:0]
	}, 1)
}

func BenchmarkWrap(b *testing.B) {
	err := errbase.New("error")
	for b.Loop() {
		_ = Wrap(err)
	}
}

func BenchmarkEnsure(b *testing.B) {
	err := errbase.New("error")
	err = Ensure(err)
	for b.Loop() {
		_ = Ensure(err)
	}
}

func BenchmarkFrames(b *testing.B) {
	err := errbase.New("error")
	err = Wrap(err)
	for b.Loop() {
		_ = Frames(err)
	}
}

func BenchmarkVerbose(b *testing.B) {
	err := errbase.New("error")
	err = Wrap(err)
	v, _ := assert.ErrorAsType[errverbose.AppendInterface](b, err)
	var buf []byte
	for b.Loop() {
		buf = v.ErrorVerboseAppend(buf)
		buf = buf[:0]
	}
}
