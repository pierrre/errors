package errors_test

import (
	"errors"
	"fmt"
	"io/fs"
	"testing"

	. "github.com/pierrre/errors"
	"github.com/pierrre/errors/errbase"
	"github.com/pierrre/errors/errstack"
)

func TestNew(t *testing.T) {
	err := New("error")
	s := err.Error()
	expected := "error"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
	sfs := errstack.Frames(err)
	if len(sfs) != 1 {
		t.Fatalf("unexpected length: got %d, want %d", len(sfs), 1)
	}
}

func TestNewf(t *testing.T) {
	err := Newf("error %d", 1)
	s := err.Error()
	expected := "error 1"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
	sfs := errstack.Frames(err)
	if len(sfs) != 1 {
		t.Fatalf("unexpected length: got %d, want %d", len(sfs), 1)
	}
}

func ExampleNew() {
	err := errors.New("error")
	fmt.Println(err)
	// Output: error
}

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

func TestWrapf(t *testing.T) {
	err := errbase.New("error")
	err = Wrapf(err, "test %d", 1)
	s := err.Error()
	expected := "test 1: error"
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

func TestAs(t *testing.T) {
	err := errbase.New("error")
	err = &fs.PathError{Err: err}
	err = Wrap(err, "test")
	var pathError *fs.PathError
	ok := As(err, &pathError)
	if !ok {
		t.Fatal("not ok")
	}
}

func TestIs(t *testing.T) {
	errBase := errbase.New("error")
	err := Wrap(errBase, "test")
	ok := Is(err, errBase)
	if !ok {
		t.Fatal("not ok")
	}
}

func TestUnwrap(t *testing.T) {
	errBase := errbase.New("error")
	err := Wrap(errBase, "test")
	err = Unwrap(err)
	err = Unwrap(err)
	if err != errBase { //nolint:errorlint // We want to compare this error.
		t.Fatal("not equal")
	}
}
