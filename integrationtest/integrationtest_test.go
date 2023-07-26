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
	"github.com/pierrre/errors/internal/errtest"
)

func init() {
	errtest.Configure()
}

var errTest error

func init() {
	err := errors.Join(
		errors.New("error a"),
		errors.New("error b"),
	)
	err = errors.Wrap(err, "test")
	err = errignore.Wrap(err)
	err = errtmp.Wrap(err, true)
	err = errtag.Wrap(err, "a", "b")
	err = errval.Wrap(err, "c", "d")
	errTest = err
}

func TestError(t *testing.T) {
	assert.ErrorEqual(t, errTest, "test: error a\nerror b")
}

func TestVerbose(t *testing.T) {
	s := errverbose.String(errTest)
	assert.RegexpMatch(t, `^test: error a\nerror b\nvalue c = d\ntag a = b\ntemporary = true\nignored\n\nSub error 0: error a\nstack\n(\t.+ .+:\d+\n)+\n\nSub error 1: error b\nstack\n(\t.+ .+:\d+\n)+\n$`, s)
}

func TestStack(t *testing.T) {
	sfs := errstack.Frames(errTest)
	assert.SliceLen(t, sfs, 2)
	for _, sf := range sfs {
		f, _ := sf.Next()
		assert.Equal(t, f.Function, "github.com/pierrre/errors/integrationtest.init.1")
	}
}

func TestIgnore(t *testing.T) {
	assert.True(t, errignore.Is(errTest))
}

func TestTemporary(t *testing.T) {
	assert.True(t, errtmp.Is(errTest))
}

func TestTag(t *testing.T) {
	tags := errtag.Get(errTest)
	assert.MapEqual(t, tags, map[string]string{"a": "b"})
}

func TestValue(t *testing.T) {
	values := errval.Get(errTest)
	assert.MapEqual(t, values, map[string]any{"c": "d"})
}

func TestErrorAllocs(t *testing.T) {
	var res string
	assert.AllocsPerRun(t, 100, func() {
		res = errTest.Error()
	}, 0)
	runtime.KeepAlive(res)
}

func TestVerboseAllocs(t *testing.T) {
	assert.AllocsPerRun(t, 100, func() {
		errverbose.Write(io.Discard, errTest)
	}, 7)
}

func TestStackAllocs(t *testing.T) {
	var res []*runtime.Frames
	assert.AllocsPerRun(t, 100, func() {
		errstack.Frames(errTest)
	}, 6)
	runtime.KeepAlive(res)
}

func TestIgnoreAllocs(t *testing.T) {
	var res bool
	assert.AllocsPerRun(t, 100, func() {
		res = errignore.Is(errTest)
	}, 1)
	runtime.KeepAlive(res)
}

func TestTemporaryAllocs(t *testing.T) {
	var res bool
	assert.AllocsPerRun(t, 100, func() {
		res = errtmp.Is(errTest)
	}, 1)
	runtime.KeepAlive(res)
}

func TestTagAllocs(t *testing.T) {
	var res map[string]string
	assert.AllocsPerRun(t, 100, func() {
		res = errtag.Get(errTest)
	}, 2)
	runtime.KeepAlive(res)
}

func TestValueAllocs(t *testing.T) {
	var res map[string]any
	assert.AllocsPerRun(t, 100, func() {
		res = errval.Get(errTest)
	}, 2)
	runtime.KeepAlive(res)
}

func BenchmarkError(b *testing.B) {
	var res string
	for i := 0; i < b.N; i++ {
		res = errTest.Error()
	}
	runtime.KeepAlive(res)
}

func BenchmarkVerbose(b *testing.B) {
	for i := 0; i < b.N; i++ {
		errverbose.Write(io.Discard, errTest)
	}
}

func BenchmarkStack(b *testing.B) {
	var res []*runtime.Frames
	for i := 0; i < b.N; i++ {
		res = errstack.Frames(errTest)
	}
	runtime.KeepAlive(res)
}

func BenchmarkIgnore(b *testing.B) {
	var res bool
	for i := 0; i < b.N; i++ {
		res = errignore.Is(errTest)
	}
	runtime.KeepAlive(res)
}

func BenchmarkTemporary(b *testing.B) {
	var res bool
	for i := 0; i < b.N; i++ {
		res = errtmp.Is(errTest)
	}
	runtime.KeepAlive(res)
}

func BenchmarkTag(b *testing.B) {
	var res map[string]string
	for i := 0; i < b.N; i++ {
		res = errtag.Get(errTest)
	}
	runtime.KeepAlive(res)
}

func BenchmarkValue(b *testing.B) {
	var res map[string]any
	for i := 0; i < b.N; i++ {
		res = errval.Get(errTest)
	}
	runtime.KeepAlive(res)
}
