package errtmp_test

import (
	"fmt"
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
	// TODO: use assert.Equal with Go 1.20.
	if err2 != err1 { //nolint:errorlint // We want to compare the error.
		t.Fatal("error not equal")
	}
}

func Example() {
	err := errbase.New("error")
	err = Wrap(err, true)
	temporary := Is(err)
	fmt.Println(temporary)
	// Output: true
}
