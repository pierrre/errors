package errors

import (
	"io"
	"testing"
)

func TestVerbose(t *testing.T) {
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

type testVerbose struct {
	error
}

func (v *testVerbose) ErrorVerbose(w io.Writer) {
	_, _ = io.WriteString(w, "verbose\n")
}
