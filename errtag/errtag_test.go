package errtag_test

import (
	"fmt"
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

func Test(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err, "foo", "bar")
	tags := Get(err)
	assert.MapEqual(t, tags, map[string]string{
		"foo": "bar",
	})
}

func Example() {
	err := errbase.New("error")
	err = Wrap(err, "foo", "bar")
	tags := Get(err)
	fmt.Println(tags["foo"])
	// Output: bar
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
