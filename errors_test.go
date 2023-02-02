package errors_test

import (
	"fmt"
	"io/fs"
	"testing"

	"github.com/pierrre/assert"
	. "github.com/pierrre/errors"
	"github.com/pierrre/errors/errbase"
	"github.com/pierrre/errors/errstack"
	"github.com/pierrre/errors/internal/errtest"
)

func init() {
	errtest.Configure()
}

func TestNew(t *testing.T) {
	err := New("error")
	assert.ErrorEqual(t, err, "error")
	sfs := errstack.Frames(err)
	assert.SliceLen(t, sfs, 1)
}

func TestNewf(t *testing.T) {
	err := Newf("error %d", 1)
	assert.ErrorEqual(t, err, "error 1")
	sfs := errstack.Frames(err)
	assert.SliceLen(t, sfs, 1)
}

func ExampleNew() {
	err := New("error")
	fmt.Println(err)
	// Output: error
}

func TestWrap(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err, "test")
	assert.ErrorEqual(t, err, "test: error")
	sfs := errstack.Frames(err)
	assert.SliceLen(t, sfs, 1)
}

func TestWrapf(t *testing.T) {
	err := errbase.New("error")
	err = Wrapf(err, "test %d", 1)
	assert.ErrorEqual(t, err, "test 1: error")
	sfs := errstack.Frames(err)
	assert.SliceLen(t, sfs, 1)
}

func ExampleWrap() {
	err := errbase.New("error")
	err = Wrap(err, "wrap")
	fmt.Println(err)
	// Output: wrap: error
}

func TestAs(t *testing.T) {
	err := errbase.New("error")
	err = &fs.PathError{Err: err}
	err = Wrap(err, "test")
	var pathError *fs.PathError
	ok := As(err, &pathError)
	assert.True(t, ok)
}

func TestIs(t *testing.T) {
	errBase := errbase.New("error")
	err := Wrap(errBase, "test")
	ok := Is(err, errBase)
	assert.True(t, ok)
}

func TestJoin(t *testing.T) {
	err := Join(New("error 1"), New("error 2"))
	errUnwrap, _ := assert.Type[interface {
		Unwrap() []error
	}](t, err)
	errs := errUnwrap.Unwrap()
	assert.SliceLen(t, errs, 2)
}

func TestUnwrap(t *testing.T) {
	errBase := errbase.New("error")
	err := Wrap(errBase, "test")
	err = Unwrap(err)
	err = Unwrap(err)
	assert.Equal(t, err, errBase)
}
