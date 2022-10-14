package errors

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestVerbose(t *testing.T) {
	err := newBase("error")
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

func TestVerboseString(t *testing.T) {
	err := newBase("error")
	err = &testVerbose{
		error: err,
	}
	s := VerboseString(err)
	expected := "error\nverbose\n"
	if s != expected {
		t.Fatalf("unexpected verbose message:\ngot: %q\nwant: %q", s, expected)
	}
}

func TestVerboseFormatter(t *testing.T) {
	err := newBase("error")
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

type testVerbose struct {
	error
}

func (v *testVerbose) ErrorVerbose(w io.Writer) {
	_, _ = io.WriteString(w, "verbose\n")
}
