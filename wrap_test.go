package errors_test

import (
	"fmt"
	"testing"

	. "github.com/pierrre/errors"
	"github.com/pierrre/errors/errbase"
	"github.com/pierrre/errors/errstack"
)

func TestWrap(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err, "test")
	s := err.Error()
	expected := "test: error"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
	sfs := errstack.Frames(err)
	if len(sfs) != 1 {
		t.Fatalf("unexpected length: got %d, want %d", len(sfs), 1)
	}
}

func ExampleWrap() {
	err := errbase.New("error")
	err = Wrap(err, "wrap")
	fmt.Println(err)
	// Output: wrap: error
}
