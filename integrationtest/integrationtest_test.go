package integrationtest

import (
	"io"
	"runtime"
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/errors"
	"github.com/pierrre/errors/errignore"
	"github.com/pierrre/errors/errstack"
	"github.com/pierrre/errors/errtag"
	"github.com/pierrre/errors/errtmp"
	"github.com/pierrre/errors/errval"
	"github.com/pierrre/errors/errverbose"
)

func newTestError() error {
	err := errors.Join(
		errors.New("error a"),
		errors.New("error b"),
	)
	err = errors.Wrap(err, "test")
	err = errignore.Wrap(err)
	err = errtmp.Wrap(err, true)
	err = errtag.Wrap(err, "a", "b")
	err = errval.Wrap(err, "c", "d")
	return err
}

func TestError(t *testing.T) {
	err := newTestError()
	assert.ErrorEqual(t, err, "test: error a\nerror b")
}

func TestVerbose(t *testing.T) {
	err := newTestError()
	s := errverbose.String(err)
	assert.RegexpMatch(t, `^test: error a\nerror b\nvalue c = \(string\) \(len=1\) "d"\ntag a = b\ntemporary = true\nignored\nstack\n(\t.+ .+:\d+\n)+\n\nSub error 0: error a\nstack\n(\t.+ .+:\d+\n)+\n\nSub error 1: error b\nstack\n(\t.+ .+:\d+\n)+\n$`, s)
}

func TestStack(t *testing.T) {
	err := newTestError()
	sfs := errstack.Frames(err)
	assert.SliceLen(t, sfs, 3)
	for _, sf := range sfs {
		f, _ := sf.Next()
		assert.Equal(t, f.Function, "github.com/pierrre/errors/integrationtest.newTestError")
	}
}

func TestIgnore(t *testing.T) {
	err := newTestError()
	assert.True(t, errignore.Is(err))
}

func TestTemporary(t *testing.T) {
	err := newTestError()
	assert.True(t, errtmp.Is(err))
}

func TestTag(t *testing.T) {
	err := newTestError()
	tags := errtag.Get(err)
	assert.MapEqual(t, tags, map[string]string{"a": "b"})
}

func TestValue(t *testing.T) {
	err := newTestError()
	values := errval.Get(err)
	assert.MapEqual(t, values, map[string]any{"c": "d"})
}

func TestNewAllocs(t *testing.T) {
	var res error
	assert.AllocsPerRun(t, 100, func() {
		res = newTestError()
	}, 18)
	runtime.KeepAlive(res)
}

func TestErrorAllocs(t *testing.T) {
	err := newTestError()
	var res string
	assert.AllocsPerRun(t, 100, func() {
		res = err.Error()
	}, 0)
	runtime.KeepAlive(res)
}

func TestVerboseAllocs(t *testing.T) {
	err := newTestError()
	assert.AllocsPerRun(t, 100, func() {
		errverbose.Write(io.Discard, err)
	}, 3)
}

func TestStackAllocs(t *testing.T) {
	err := newTestError()
	var res []*runtime.Frames
	assert.AllocsPerRun(t, 100, func() {
		errstack.Frames(err)
	}, 6)
	runtime.KeepAlive(res)
}

func TestIgnoreAllocs(t *testing.T) {
	err := newTestError()
	var res bool
	assert.AllocsPerRun(t, 100, func() {
		res = errignore.Is(err)
	}, 1)
	runtime.KeepAlive(res)
}

func TestTemporaryAllocs(t *testing.T) {
	err := newTestError()
	var res bool
	assert.AllocsPerRun(t, 100, func() {
		res = errtmp.Is(err)
	}, 1)
	runtime.KeepAlive(res)
}

func TestTagAllocs(t *testing.T) {
	err := newTestError()
	var res map[string]string
	assert.AllocsPerRun(t, 100, func() {
		res = errtag.Get(err)
	}, 2)
	runtime.KeepAlive(res)
}

func TestValueAllocs(t *testing.T) {
	err := newTestError()
	var res map[string]any
	assert.AllocsPerRun(t, 100, func() {
		res = errval.Get(err)
	}, 2)
	runtime.KeepAlive(res)
}

func BenchmarkNew(b *testing.B) {
	var res error
	for i := 0; i < b.N; i++ {
		res = newTestError()
	}
	runtime.KeepAlive(res)
}

func BenchmarkError(b *testing.B) {
	err := newTestError()
	var res string
	for i := 0; i < b.N; i++ {
		res = err.Error()
	}
	runtime.KeepAlive(res)
}

func BenchmarkVerbose(b *testing.B) {
	err := newTestError()
	for i := 0; i < b.N; i++ {
		errverbose.Write(io.Discard, err)
	}
}

func BenchmarkStack(b *testing.B) {
	err := newTestError()
	var res []*runtime.Frames
	for i := 0; i < b.N; i++ {
		res = errstack.Frames(err)
	}
	runtime.KeepAlive(res)
}

func BenchmarkIgnore(b *testing.B) {
	err := newTestError()
	var res bool
	for i := 0; i < b.N; i++ {
		res = errignore.Is(err)
	}
	runtime.KeepAlive(res)
}

func BenchmarkTemporary(b *testing.B) {
	err := newTestError()
	var res bool
	for i := 0; i < b.N; i++ {
		res = errtmp.Is(err)
	}
	runtime.KeepAlive(res)
}

func BenchmarkTag(b *testing.B) {
	err := newTestError()
	var res map[string]string
	for i := 0; i < b.N; i++ {
		res = errtag.Get(err)
	}
	runtime.KeepAlive(res)
}

func BenchmarkValue(b *testing.B) {
	err := newTestError()
	var res map[string]any
	for i := 0; i < b.N; i++ {
		res = errval.Get(err)
	}
	runtime.KeepAlive(res)
}
