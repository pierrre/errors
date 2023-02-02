package errignore_test

import (
	"fmt"
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/errors"
	"github.com/pierrre/errors/errbase"
	. "github.com/pierrre/errors/errignore"
	"github.com/pierrre/errors/errverbose"
	"github.com/pierrre/errors/internal/errtest"
)

func init() {
	errtest.Configure()
}

func Test(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err)
	ignored := Is(err)
	assert.True(t, ignored)
}

func TestNil(t *testing.T) {
	err := Wrap(nil)
	assert.NoError(t, err)
}

func TestFalse(t *testing.T) {
	err := errbase.New("error")
	ignored := Is(err)
	assert.False(t, ignored)
}

func TestError(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err)
	assert.ErrorEqual(t, err, "error")
}

func TestVerbose(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err)
	var v errverbose.Interface
	assert.ErrorAs(t, err, &v)
	s := v.ErrorVerbose()
	assert.Equal(t, s, "ignored")
}

func TestUnwrap(t *testing.T) {
	err1 := errbase.New("error")
	err2 := Wrap(err1)
	err2 = errors.Unwrap(err2)
	assert.Equal(t, err2, err1)
}

func Example() {
	err := errbase.New("error")
	err = Wrap(err)
	ignored := Is(err)
	fmt.Println(ignored)
	// Output: true
}
