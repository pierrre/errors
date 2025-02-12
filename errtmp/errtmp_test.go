package errtmp_test

import (
	"fmt"
	"io"
	"runtime"
	"strings"
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/errors"
	"github.com/pierrre/errors/errbase"
	. "github.com/pierrre/errors/errtmp"
	"github.com/pierrre/errors/errverbose"
)

func Example() {
	err := errbase.New("error")
	err = Wrap(err, true)
	temporary := Is(err)
	fmt.Println(temporary)
	// Output: true
}

func TestTrue(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err, true)
	temporary := Is(err)
	assert.True(t, temporary)
}

func TestFalse(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err, false)
	temporary := Is(err)
	assert.False(t, temporary)
}

func TestDefault(t *testing.T) {
	err := errbase.New("error")
	temporary := Is(err)
	assert.True(t, temporary)
}

func TestNil(t *testing.T) {
	err := Wrap(nil, true)
	assert.NoError(t, err)
}

func TestError(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err, true)
	assert.ErrorEqual(t, err, "error")
}

func TestVerbose(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err, true)
	var v errverbose.Interface
	assert.ErrorAs(t, err, &v)
	sb := new(strings.Builder)
	v.ErrorVerbose(sb)
	s := sb.String()
	assert.Equal(t, s, "temporary = true")
}

func TestUnwrap(t *testing.T) {
	err1 := errbase.New("error")
	err2 := Wrap(err1, true)
	err2 = errors.Unwrap(err2)
	assert.Equal(t, err2, err1)
}

func TestWrapAllocs(t *testing.T) {
	err := errbase.New("error")
	var res error
	assert.AllocsPerRun(t, 100, func() {
		res = Wrap(err, true)
	}, 1)
	runtime.KeepAlive(res)
}

func TestIsAllocs(t *testing.T) {
	err := errbase.New("error")
	var res bool
	assert.AllocsPerRun(t, 100, func() {
		res = Is(err)
	}, 1)
	runtime.KeepAlive(res)
}

func TestVerboseAllocs(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err, true)
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
		res = Wrap(err, true)
	}
	runtime.KeepAlive(res)
}

func BenchmarkIs(b *testing.B) {
	err := errbase.New("error")
	err = Wrap(err, true)
	var res bool
	for b.Loop() {
		res = Is(err)
	}
	runtime.KeepAlive(res)
}

func BenchmarkVerbose(b *testing.B) {
	err := errbase.New("error")
	err = Wrap(err, true)
	var v errverbose.Interface
	assert.ErrorAs(b, err, &v)
	for b.Loop() {
		v.ErrorVerbose(io.Discard)
	}
}
