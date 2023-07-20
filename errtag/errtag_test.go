package errtag_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/errors"
	"github.com/pierrre/errors/errbase"
	. "github.com/pierrre/errors/errtag"
	"github.com/pierrre/errors/errverbose"
	"github.com/pierrre/errors/internal/errtest"
)

func init() {
	errtest.Configure()
}

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
	s := v.ErrorVerbose()
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
	var res string
	assert.AllocsPerRun(t, 100, func() {
		res = v.ErrorVerbose()
	}, 1)
	runtime.KeepAlive(res)
}

func BenchmarkWrap(b *testing.B) {
	err := errbase.New("error")
	var res error
	for i := 0; i < b.N; i++ {
		res = Wrap(err, "foo", "bar")
	}
	runtime.KeepAlive(res)
}

func BenchmarkWrapInt(b *testing.B) {
	err := errbase.New("error")
	var res error
	for i := 0; i < b.N; i++ {
		res = WrapInt(err, "foo", 123)
	}
	runtime.KeepAlive(res)
}

func BenchmarkWrapInt64(b *testing.B) {
	err := errbase.New("error")
	var res error
	for i := 0; i < b.N; i++ {
		res = WrapInt64(err, "foo", 123)
	}
	runtime.KeepAlive(res)
}

func BenchmarkWrapFloat64(b *testing.B) {
	err := errbase.New("error")
	var res error
	for i := 0; i < b.N; i++ {
		res = WrapFloat64(err, "foo", 12.3)
	}
	runtime.KeepAlive(res)
}

func BenchmarkWrapBool(b *testing.B) {
	err := errbase.New("error")
	var res error
	for i := 0; i < b.N; i++ {
		res = WrapBool(err, "foo", true)
	}
	runtime.KeepAlive(res)
}

func BenchmarkGet(b *testing.B) {
	err := errbase.New("error")
	err = Wrap(err, "foo", "bar")
	var res map[string]string
	for i := 0; i < b.N; i++ {
		res = Get(err)
	}
	runtime.KeepAlive(res)
}

func BenchmarkVerbose(b *testing.B) {
	err := errbase.New("error")
	err = Wrap(err, "foo", "bar")
	var v errverbose.Interface
	assert.ErrorAs(b, err, &v)
	var res string
	for i := 0; i < b.N; i++ {
		res = v.ErrorVerbose()
	}
	runtime.KeepAlive(res)
}
