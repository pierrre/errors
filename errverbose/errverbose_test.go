package errverbose

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pierrre/errors/errbase"
)

func TestWrite(t *testing.T) {
	err := errbase.New("error")
	err = &testVerbose{
		error: err,
	}
	buf := new(strings.Builder)
	Write(buf, err)
	s := buf.String()
	expected := "error\nverbose\n"
	if s != expected {
		t.Fatalf("unexpected verbose message:\ngot: %q\nwant: %q", s, expected)
	}
}

//nolint:testableexamples // The output contains a stack trace, which is not stable.
func ExampleWrite() {
	err := errbase.New("error")
	buf := new(strings.Builder)
	Write(buf, err)
	s := buf.String()
	fmt.Println(s)
}

func TestString(t *testing.T) {
	err := errbase.New("error")
	err = &testVerbose{
		error: err,
	}
	s := String(err)
	expected := "error\nverbose\n"
	if s != expected {
		t.Fatalf("unexpected verbose message:\ngot: %q\nwant: %q", s, expected)
	}
}

//nolint:testableexamples // The output contains a stack trace, which is not stable.
func ExampleString() {
	err := errbase.New("error")
	s := String(err)
	fmt.Println(s)
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
	expected := "error\nverbose\n"
	if s != expected {
		t.Fatalf("unexpected verbose message:\ngot: %q\nwant: %q", s, expected)
	}
}

//nolint:testableexamples // The output contains a stack trace, which is not stable.
func ExampleFormatter() {
	err := errbase.New("error")
	f := Formatter(err)
	fmt.Println(f)
}

type testVerbose struct {
	error
}

func (v *testVerbose) ErrorVerbose() string {
	return "verbose"
}
