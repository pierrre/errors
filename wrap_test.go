package errors

import (
	"fmt"
	"testing"
)

func TestWrap(t *testing.T) {
	err := newBase("error")
	err = Wrap(err, "test")
	s := err.Error()
	expected := "test: error"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
	sfs := StackFrames(err)
	if len(sfs) != 1 {
		t.Fatalf("unexpected length: got %d, want %d", len(sfs), 1)
	}
}

func ExampleWrap() {
	err := New("error")
	err = Wrap(err, "wrap")
	fmt.Println(err)
	// Output: wrap: error
}
