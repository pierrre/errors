package errverbose_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/errors/errbase"
	. "github.com/pierrre/errors/errverbose"
	"github.com/pierrre/errors/internal/errtest"
)

func init() {
	errtest.Configure()
}

func TestWrite(t *testing.T) {
	err := errbase.New("error")
	err = &testVerbose{
		error: err,
	}
	buf := new(strings.Builder)
	Write(buf, err)
	s := buf.String()
	assert.Equal(t, s, "error\nverbose\n")
}

func ExampleWrite() {
	err := errbase.New("error")
	buf := new(strings.Builder)
	Write(buf, err)
	s := buf.String()
	fmt.Println(s)
	// Output: error
}

func TestString(t *testing.T) {
	err := errbase.New("error")
	err = &testVerbose{
		error: err,
	}
	s := String(err)
	assert.Equal(t, s, "error\nverbose\n")
}

func ExampleString() {
	err := errbase.New("error")
	s := String(err)
	fmt.Println(s)
	// Output: error
}

func TestFormatter(t *testing.T) {
	err := errbase.New("error")
	err = &testVerbose{
		error: err,
	}
	f := Formatter(err)
	buf := new(strings.Builder)
	_, _ = fmt.Fprintf(buf, "%v", f)
	s := buf.String()
	assert.Equal(t, s, "error\nverbose\n")
}

func ExampleFormatter() {
	err := errbase.New("error")
	f := Formatter(err)
	fmt.Println(f)
	// Output: error
}

type testVerbose struct {
	error
}

func (v *testVerbose) ErrorVerbose() string {
	return "verbose"
}
