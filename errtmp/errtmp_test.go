package errtmp_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/errors"
	"github.com/pierrre/errors/errbase"
	. "github.com/pierrre/errors/errtmp"
	"github.com/pierrre/errors/errverbose"
	"github.com/pierrre/errors/internal/errtest"
)

func init() {
	errtest.Configure()
}

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
	s := v.ErrorVerbose()
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
		res = Wrap(err, true)
	}
	runtime.KeepAlive(res)
}

func BenchmarkIs(b *testing.B) {
	err := errbase.New("error")
	err = Wrap(err, true)
	var res bool
	for i := 0; i < b.N; i++ {
		res = Is(err)
	}
	runtime.KeepAlive(res)
}

func BenchmarkVerbose(b *testing.B) {
	err := errbase.New("error")
	err = Wrap(err, true)
	var v errverbose.Interface
	assert.ErrorAs(b, err, &v)
	var res string
	for i := 0; i < b.N; i++ {
		res = v.ErrorVerbose()
	}
	runtime.KeepAlive(res)
}
