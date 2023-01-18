package errignore_test

import (
	"fmt"
	"testing"

	"github.com/pierrre/errors"
	"github.com/pierrre/errors/errbase"
	. "github.com/pierrre/errors/errignore"
	"github.com/pierrre/errors/errverbose"
)

func Test(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err)
	ignored := Is(err)
	if !ignored {
		t.Fatalf("unexpected ignored: got %t, want %t", ignored, true)
	}
}

func TestNil(t *testing.T) {
	err := Wrap(nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFalse(t *testing.T) {
	err := errbase.New("error")
	ignored := Is(err)
	if ignored {
		t.Fatalf("unexpected ignored: got %t, want %t", ignored, false)
	}
}

func TestError(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err)
	s := err.Error()
	expected := "error"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
}

func TestVerbose(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err)
	var v errverbose.Interface
	ok := errors.As(err, &v)
	if !ok {
		t.Fatal("not a Verbose")
	}
	s := v.ErrorVerbose()
	expected := "ignored"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
}

func TestUnwrap(t *testing.T) {
	err1 := errbase.New("error")
	err2 := Wrap(err1)
	err2 = errors.Unwrap(err2)
	if err2 != err1 { //nolint:errorlint // We want to compare the error.
		t.Fatal("error not equal")
	}
}

func Example() {
	err := errbase.New("error")
	err = Wrap(err)
	ignored := Is(err)
	fmt.Println(ignored)
	// Output: true
}
