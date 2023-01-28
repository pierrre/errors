package errstack_test

import (
	"fmt"
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/errors"
	"github.com/pierrre/errors/errbase"
	. "github.com/pierrre/errors/errstack"
	"github.com/pierrre/errors/errverbose"
	"github.com/pierrre/errors/internal/errtest"
)

func init() {
	errtest.Configure()
}

func Test(t *testing.T) {
	err := errbase.New("error")
	err = Ensure(err)
	err = Ensure(err)
	sfs := Frames(err)
	assert.SliceLen(t, sfs, 1)
	sf := sfs[0]
	assert.NotZero(t, sf)
	f, _ := sf.Next()
	assert.Equal(t, f.Function, "github.com/pierrre/errors/errstack_test.Test")
}

func TestNil(t *testing.T) {
	err := Wrap(nil)
	assert.NoError(t, err)
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
	assert.RegexpMatch(t, `^stack\n(\t.+ .+:\d+\n)+$`, s)
}

func TestStackFrames(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err)
	var sErr interface {
		StackFrames() []uintptr
	}
	assert.ErrorAs(t, err, &sErr)
	pcs := sErr.StackFrames()
	assert.SliceNotEmpty(t, pcs)
}

func Example() {
	err := errors.New("error")
	err = Wrap(err)
	fmt.Println(err)
	sfs := Frames(err)
	fmt.Println(len(sfs))
	// Output:
	// error
	// 2
}
