package errors

import (
	"testing"
)

func TestWrap(t *testing.T) {
	err := newBase("error")
	err = Wrap(err, "test1")
	err = Wrapf(err, "%s", "test2")
	s := err.Error()
	expected := "test2: test1: error"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
	sfs := StackFrames(err)
	if len(sfs) != 1 {
		t.Fatalf("unexpected length: got %d, want %d", len(sfs), 1)
	}
}
