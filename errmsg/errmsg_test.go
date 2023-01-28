package errmsg_test

import (
	"fmt"
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/errors"
	"github.com/pierrre/errors/errbase"
	. "github.com/pierrre/errors/errmsg"
	"github.com/pierrre/errors/internal/errtest"
)

func init() {
	errtest.Configure()
}

func Test(t *testing.T) {
	err := errbase.New("error")
	err = Wrapf(err, "test %d", 1)
	assert.ErrorEqual(t, err, "test 1: error")
}

func TestNil(t *testing.T) {
	err := Wrap(nil, "test")
	assert.NoError(t, err)
}

func TestEmpty(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err, "")
	assert.ErrorEqual(t, err, "error")
}

func TestUnwrap(t *testing.T) {
	err1 := errbase.New("error")
	err2 := Wrap(err1, "test")
	err2 = errors.Unwrap(err2)
	// TODO: use assert.Equal with Go 1.20.
	if err2 != err1 { //nolint:errorlint // We want to compare the error.
		t.Fatal("error not equal")
	}
}

func Example() {
	err := errbase.New("error")
	err = Wrap(err, "message")
	fmt.Println(err)
	// Output: message: error
}
