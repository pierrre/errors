package errtag_test

import (
	"fmt"
	"io"
	"runtime"
	"strings"
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/errors"
	"github.com/pierrre/errors/errbase"
	. "github.com/pierrre/errors/errtag"
	"github.com/pierrre/errors/errverbose"
)

func Example() {
	err := errbase.New("error")
	err = Wrap(err, "foo", "bar")
	tags := Get(err)
	fmt.Println(tags["foo"])
	// Output: bar
}

func Test(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err, "foo", "bar")
	tags := Get(err)
	assert.MapEqual(t, tags, map[string]string{
		"foo": "bar",
	})
}

func TestInt(t *testing.T) {
	err := errbase.New("error")
	err = WrapInt(err, "foo", 123)
	tags := Get(err)
	assert.MapEqual(t, tags, map[string]string{
		"foo": "123",
	})
}

func TestInt64(t *testing.T) {
	err := errbase.New("error")
	err = WrapInt64(err, "foo", 123)
	tags := Get(err)
	assert.MapEqual(t, tags, map[string]string{
		"foo": "123",
	})
}

func TestFloat64(t *testing.T) {
	err := errbase.New("error")
	err = WrapFloat64(err, "foo", 12.3)
	tags := Get(err)
	assert.MapEqual(t, tags, map[string]string{
		"foo": "12.3",
	})
}

func TestBool(t *testing.T) {
	err := errbase.New("error")
	err = WrapBool(err, "foo", true)
	tags := Get(err)
	assert.MapEqual(t, tags, map[string]string{
		"foo": "true",
	})
}

func TestOverWrite(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err, "test", "1")
	err = Wrap(err, "test", "2")
	tags := Get(err)
	assert.MapEqual(t, tags, map[string]string{
		"test": "2",
	})
}

func TestNil(t *testing.T) {
	err := Wrap(nil, "foo", "bar")
	assert.NoError(t, err)
}

func TestEmpty(t *testing.T) {
	err := errbase.New("error")
	tags := Get(err)
	assert.MapEmpty(t, tags)
}

func TestError(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err, "foo", "bar")
	assert.ErrorEqual(t, err, "error")
}

func TestVerbose(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err, "foo", "bar")
	var v errverbose.Interface
	assert.ErrorAs(t, err, &v)
	sb := new(strings.Builder)
	v.ErrorVerbose(sb)
	s := sb.String()
	assert.Equal(t, s, "tag foo = bar")
}

func TestJoin(t *testing.T) {
	err := Wrap(
		errors.Join(
			Wrap(
				errors.New("error"),
				"foo",
				"baz",
			),
			Wrap(
				errors.New("error"),
				"aaa",
				"bbb",
			),
		),
		"foo",
		"bar",
	)
	tags := Get(err)
	assert.MapEqual(t, tags, map[string]string{
		"foo": "bar",
		"aaa": "bbb",
	})
}

func TestWrapAllocs(t *testing.T) {
	err := errbase.New("error")
	var res error
	assert.AllocsPerRun(t, 100, func() {
		res = Wrap(err, "foo", "bar")
	}, 1)
	runtime.KeepAlive(res)
}

func TestWrapIntAllocs(t *testing.T) {
	err := errbase.New("error")
	var res error
	assert.AllocsPerRun(t, 100, func() {
		res = WrapInt(err, "foo", 123)
	}, 2)
	runtime.KeepAlive(res)
}

func TestWrapInt64Allocs(t *testing.T) {
	err := errbase.New("error")
	var res error
	assert.AllocsPerRun(t, 100, func() {
		res = WrapInt64(err, "foo", 123)
	}, 2)
	runtime.KeepAlive(res)
}

func TestWrapFloat64Allocs(t *testing.T) {
	err := errbase.New("error")
	var res error
	assert.AllocsPerRun(t, 100, func() {
		res = WrapFloat64(err, "foo", 12.3)
	}, 3)
	runtime.KeepAlive(res)
}

func TestWrapBoolAllocs(t *testing.T) {
	err := errbase.New("error")
	var res error
	assert.AllocsPerRun(t, 100, func() {
		res = WrapBool(err, "foo", true)
	}, 1)
	runtime.KeepAlive(res)
}

func TestGetAllocs(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err, "foo", "bar")
	var res map[string]string
	assert.AllocsPerRun(t, 100, func() {
		res = Get(err)
	}, 2)
	runtime.KeepAlive(res)
}

func TestVerboseAllocs(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err, "foo", "bar")
	var v errverbose.Interface
	assert.ErrorAs(t, err, &v)
	assert.AllocsPerRun(t, 100, func() {
		v.ErrorVerbose(io.Discard)
	}, 0)
}

func BenchmarkWrap(b *testing.B) {
	err := errbase.New("error")
	var res error
	for b.Loop() {
		res = Wrap(err, "foo", "bar")
	}
	runtime.KeepAlive(res)
}

func BenchmarkWrapInt(b *testing.B) {
	err := errbase.New("error")
	var res error
	for b.Loop() {
		res = WrapInt(err, "foo", 123)
	}
	runtime.KeepAlive(res)
}

func BenchmarkWrapInt64(b *testing.B) {
	err := errbase.New("error")
	var res error
	for b.Loop() {
		res = WrapInt64(err, "foo", 123)
	}
	runtime.KeepAlive(res)
}

func BenchmarkWrapFloat64(b *testing.B) {
	err := errbase.New("error")
	var res error
	for b.Loop() {
		res = WrapFloat64(err, "foo", 12.3)
	}
	runtime.KeepAlive(res)
}

func BenchmarkWrapBool(b *testing.B) {
	err := errbase.New("error")
	var res error
	for b.Loop() {
		res = WrapBool(err, "foo", true)
	}
	runtime.KeepAlive(res)
}

func BenchmarkGet(b *testing.B) {
	err := errbase.New("error")
	err = Wrap(err, "foo", "bar")
	var res map[string]string
	for b.Loop() {
		res = Get(err)
	}
	runtime.KeepAlive(res)
}

func BenchmarkVerbose(b *testing.B) {
	err := errbase.New("error")
	err = Wrap(err, "foo", "bar")
	var v errverbose.Interface
	assert.ErrorAs(b, err, &v)
	for b.Loop() {
		v.ErrorVerbose(io.Discard)
	}
}
