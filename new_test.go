package errors_test

import (
	"fmt"
	"testing"

	"github.com/pierrre/errors"
)

func TestNew(t *testing.T) {
	err := errors.New("error")
	s := err.Error()
	expected := "error"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
	sfs := errors.StackFrames(err)
	if len(sfs) != 1 {
		t.Fatalf("unexpected length: got %d, want %d", len(sfs), 1)
	}
}

func ExampleNew() {
	err := errors.New("error")
	fmt.Println(err)
	// Output: error
}
