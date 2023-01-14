package integrationtest

import (
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

func Test(t *testing.T) {
	err := errors.Join(
		errors.New("error a"),
		errors.New("error b"),
	)
	err = errors.Wrap(err, "test")
	err = errignore.Wrap(err)
	err = errtmp.Wrap(err, true)
	err = errtag.Wrap(err, "a", "b")
	err = errval.Wrap(err, "c", "d")
	for _, tc := range []struct {
		name string
		f    func(*testing.T, error)
	}{
		{"Error", testError},
		{"Verbose", testVerbose},
		{"Stack", testStack},
		{"Ignore", testIgnore},
		{"Temporary", testTemporary},
		{"Tag", testTag},
		{"Value", testValue},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc.f(t, err)
		})
	}
}

func testError(t *testing.T, err error) {
	t.Helper()
	assert.ErrorEqual(t, err, "test: error a\nerror b")
}

func testVerbose(t *testing.T, err error) {
	t.Helper()
	s := errverbose.String(err)
	assert.RegexpMatch(t, `^test: error a\nerror b\nvalue c = d\ntag a = b\ntemporary = true\nignored\n\nSub error 0: error a\nstack\n(\t.+ .+:\d+\n)+\n\nSub error 1: error b\nstack\n(\t.+ .+:\d+\n)+\n$`, s)
}

func testStack(t *testing.T, err error) {
	t.Helper()
	sfs := errstack.Frames(err)
	assert.SliceLen(t, sfs, 2)
	for _, sf := range sfs {
		f, _ := sf.Next()
		assert.Equal(t, f.Function, "github.com/pierrre/errors/integrationtest.Test")
	}
}

func testIgnore(t *testing.T, err error) {
	t.Helper()
	assert.True(t, errignore.Is(err))
}

func testTemporary(t *testing.T, err error) {
	t.Helper()
	assert.True(t, errtmp.Is(err))
}

func testTag(t *testing.T, err error) {
	t.Helper()
	tags := errtag.Get(err)
	assert.MapEqual(t, tags, map[string]string{"a": "b"})
}

func testValue(t *testing.T, err error) {
	t.Helper()
	values := errval.Get(err)
	assert.MapEqual(t, values, map[string]any{"c": "d"})
}
