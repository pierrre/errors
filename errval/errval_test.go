package errval_test

import (
	"fmt"
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/errors"
	"github.com/pierrre/errors/errappend"
	"github.com/pierrre/errors/errbase"
	"github.com/pierrre/errors/errmsg"
	. "github.com/pierrre/errors/errval"
	"github.com/pierrre/errors/errverbose"
)

var testSink any

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
	var v errverbose.AppendInterface
	assert.ErrorAs(t, err, &v)
	b := v.ErrorVerboseAppend(nil)
	assert.Equal(t, string(b), `value foo = [string] (len=3) "bar"`)
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

func TestErrorAppend(t *testing.T) {
	err := errmsg.Wrap(errbase.New("error"), "msg")
	err = Wrap(err, "foo", "bar")
	var i errappend.Interface
	assert.ErrorAs(t, err, &i)
	b := errappend.Append(nil, err)
	assert.Equal(t, string(b), "msg: error")
}

func TestAllInterrupt(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err, "foo", "bar")
	count := 0
	for k, v := range All(err) {
		count++
		assert.Equal(t, k, "foo")
		assert.Equal(t, v, "bar")
		break
	}
	assert.Equal(t, count, 1)
}

func TestWrapAllocs(t *testing.T) {
	err := errbase.New("error")
	var res error
	assert.AllocsPerRun(t, 100, func() {
		res = Wrap(err, "foo", "bar")
	}, 1)
	testSink = res
}

func TestGetAllocs(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err, "foo", "bar")
	var res map[string]any
	assert.AllocsPerRun(t, 100, func() {
		res = Get(err)
	}, 2)
	testSink = res
}

func TestVerboseAllocs(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err, "foo", "bar")
	var v errverbose.AppendInterface
	assert.ErrorAs(t, err, &v)
	var b []byte
	assert.AllocsPerRun(t, 100, func() {
		b = v.ErrorVerboseAppend(b)
		b = b[:0]
	}, 0)
}

func ExampleGetValue() {
	err := errbase.New("error")
	err = Wrap(err, "foo", "bar")
	val, ok := GetValue(err, "foo")
	fmt.Println(ok, val)
	// Output: true bar
}

func TestGetValue(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err, "foo", "bar")
	val, ok := GetValue(err, "foo")
	assert.True(t, ok)
	assert.Equal(t, val, "bar")
}

func TestGetValueNotFound(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err, "foo", "bar")
	val, ok := GetValue(err, "baz")
	assert.False(t, ok)
	assert.Equal(t, val, nil)
}

func TestGetValueOverWrite(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err, "test", 1)
	err = Wrap(err, "test", 2)
	val, ok := GetValue(err, "test")
	assert.True(t, ok)
	assert.Equal(t, val, 2)
}

func TestGetValueEmpty(t *testing.T) {
	err := errbase.New("error")
	val, ok := GetValue(err, "foo")
	assert.False(t, ok)
	assert.Equal(t, val, nil)
}

func TestGetValueNil(t *testing.T) {
	val, ok := GetValue(nil, "foo")
	assert.False(t, ok)
	assert.Equal(t, val, nil)
}

func TestGetValueJoin(t *testing.T) {
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
	val, ok := GetValue(err, "foo")
	assert.True(t, ok)
	assert.Equal(t, val, "bar")
	val, ok = GetValue(err, "aaa")
	assert.True(t, ok)
	assert.Equal(t, val, "bbb")
	val, ok = GetValue(err, "missing")
	assert.False(t, ok)
	assert.Equal(t, val, nil)
}

func TestGetValueAllocs(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err, "foo", "bar")
	var res any
	var ok bool
	assert.AllocsPerRun(t, 100, func() {
		res, ok = GetValue(err, "foo")
	}, 0)
	testSink = res
	testSink = ok
}

func BenchmarkWrap(b *testing.B) {
	err := errbase.New("error")
	for b.Loop() {
		_ = Wrap(err, "foo", "bar")
	}
}

func BenchmarkGet(b *testing.B) {
	err := errbase.New("error")
	err = Wrap(err, "foo", "bar")
	for b.Loop() {
		_ = Get(err)
	}
}

func BenchmarkGetValue(b *testing.B) {
	err := errbase.New("error")
	err = Wrap(err, "foo", "bar")
	for b.Loop() {
		_, _ = GetValue(err, "foo")
	}
}

func BenchmarkVerbose(b *testing.B) {
	err := errbase.New("error")
	err = Wrap(err, "foo", "bar")
	var v errverbose.AppendInterface
	assert.ErrorAs(b, err, &v)
	var buf []byte
	for b.Loop() {
		buf = v.ErrorVerboseAppend(buf)
		buf = buf[:0]
	}
}
