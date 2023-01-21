package errmsg_test

import (
	"fmt"
	"testing"

	"github.com/pierrre/errors"
	"github.com/pierrre/errors/errbase"
	. "github.com/pierrre/errors/errmsg"
)

func Test(t *testing.T) {
	err := errbase.New("error")
	err = Wrapf(err, "test %d", 1)
	s := err.Error()
	expected := "test 1: error"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
}

func TestNil(t *testing.T) {
	err := Wrap(nil, "test")
	if err != nil {
		t.Fatal(err)
	}
}

func TestEmpty(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err, "")
	s := err.Error()
	expected := "error"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
}

func TestUnwrap(t *testing.T) {
	err1 := errbase.New("error")
	err2 := Wrap(err1, "test")
	err2 = errors.Unwrap(err2)
	if err2 != err1 { //nolint:errorlint // We want to compare the error.
		t.Fatal("error not equal")
	}
}

func Example() {
	err := errbase.New("error")
	err = Wrap(err, "message")
	fmt.Println(err)
	// Output: message: error
}
