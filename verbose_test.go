package errors

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/pierrre/errors/errbase"
)

func TestVerbose(t *testing.T) {
	err := errbase.New("error")
	err = &testVerbose{
		error: err,
	}
	buf := new(strings.Builder)
	Verbose(buf, err)
	s := buf.String()
	expected := "error\nverbose\n"
	if s != expected {
		t.Fatalf("unexpected verbose message:\ngot: %q\nwant: %q", s, expected)
	}
}

//nolint:testableexamples // The output contains a stack trace, which is not stable.
func ExampleVerbose() {
	err := New("error")
	buf := new(strings.Builder)
	Verbose(buf, err)
	s := buf.String()
	fmt.Println(s)
}

func TestVerboseString(t *testing.T) {
	err := errbase.New("error")
	err = &testVerbose{
		error: err,
	}
	s := VerboseString(err)
	expected := "error\nverbose\n"
	if s != expected {
		t.Fatalf("unexpected verbose message:\ngot: %q\nwant: %q", s, expected)
	}
}

//nolint:testableexamples // The output contains a stack trace, which is not stable.
func ExampleVerboseString() {
	err := New("error")
	s := VerboseString(err)
	fmt.Println(s)
}

func TestVerboseFormatter(t *testing.T) {
	err := errbase.New("error")
	err = &testVerbose{
		error: err,
	}
	f := VerboseFormatter(err)
	buf := new(strings.Builder)
	_, _ = fmt.Fprintf(buf, "%v", f)
	s := buf.String()
	expected := "error\nverbose\n"
	if s != expected {
		t.Fatalf("unexpected verbose message:\ngot: %q\nwant: %q", s, expected)
	}
}

//nolint:testableexamples // The output contains a stack trace, which is not stable.
func ExampleVerboseFormatter() {
	err := New("error")
	f := VerboseFormatter(err)
	fmt.Println(f)
}

type testVerbose struct {
	error
}

func (v *testVerbose) ErrorVerbose(w io.Writer) {
	_, _ = io.WriteString(w, "verbose\n")
}
