package errors

import (
	"bytes"
	"regexp"
	"testing"
)

func TestStack(t *testing.T) {
	err := newBase("error")
	err = Stack(err)
	sfs := StackFrames(err)
	if len(sfs) != 1 {
		t.Fatalf("unexpected length: got %d, want %d", len(sfs), 1)
	}
	sf := sfs[0]
	if sf == nil {
		t.Fatal("no stack frames")
	}
	f, _ := sf.Next()
	expectedFunction := "github.com/pierrre/errors.TestStack"
	if f.Function != expectedFunction {
		t.Fatalf("unexpected function: got %q, want %q", f.Function, expectedFunction)
	}
}

func TestStackNil(t *testing.T) {
	err := Stack(nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStackError(t *testing.T) {
	err := newBase("error")
	err = Stack(err)
	s := err.Error()
	expected := "error"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
}

func TestStackVerbose(t *testing.T) {
	err := newBase("error")
	err = Stack(err)
	var v Verboser
	ok := As(err, &v)
	if !ok {
		t.Fatal("not a Verbose")
	}
	buf := new(bytes.Buffer)
	v.ErrorVerbose(buf)
	s := buf.String()
	expectedRegexp := regexp.MustCompile(`^stack\n(\t.+ .+:\d+\n)+$`)
	if !expectedRegexp.MatchString(s) {
		t.Fatalf("unexpected verbose message:\ngot: %q\nwant match: %q", s, expectedRegexp)
	}
}

func TestStackPCs(t *testing.T) {
	err := newBase("error")
	err = Stack(err)
	var sErr *stack
	ok := As(err, &sErr)
	if !ok {
		t.Fatal("not a stack")
	}
	pcs := sErr.StackPCs()
	if len(pcs) == 0 {
		t.Fatal("no stack PCs")
	}
}
