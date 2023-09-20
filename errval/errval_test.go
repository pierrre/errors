package errval_test

import (
	"fmt"
	"io"
	"runtime"
	"strings"
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/errors"
	"github.com/pierrre/errors/errbase"
	. "github.com/pierrre/errors/errval"
	"github.com/pierrre/errors/errverbose"
	"github.com/pierrre/errors/internal/errtest"
)

func init() {
	errtest.Configure()
}

func Example() {
	err := errbase.New("error")
	err = Wrap(err, "foo", "bar")
	vals := Get(err)
	fmt.Println(vals["foo"])
	// Output: bar
}

func Test(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err, "foo", "bar")
	vals := Get(err)
	assert.MapEqual(t, vals, map[string]any{
		"foo": "bar",
	})
}

func TestOverWrite(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err, "test", 1)
	err = Wrap(err, "test", 2)
	vals := Get(err)
	assert.MapEqual(t, vals, map[string]any{
		"test": 2,
	})
}

func TestNil(t *testing.T) {
	err := Wrap(nil, "foo", "bar")
	assert.NoError(t, err)
}

func TestEmpty(t *testing.T) {
	err := errbase.New("error")
	vals := Get(err)
	assert.MapEmpty(t, vals)
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
	assert.Equal(t, s, "value foo = bar")
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
	assert.MapEqual(t, tags, map[string]any{
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

func TestGetAllocs(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err, "foo", "bar")
	var res map[string]any
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
	for i := 0; i < b.N; i++ {
		res = Wrap(err, "foo", "bar")
	}
	runtime.KeepAlive(res)
}

func BenchmarkGet(b *testing.B) {
	err := errbase.New("error")
	err = Wrap(err, "foo", "bar")
	var res map[string]any
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
	for i := 0; i < b.N; i++ {
		v.ErrorVerbose(io.Discard)
	}
}
