package errignore

import (
	"bytes"
	"testing"

	"github.com/pierrre/errors"
)

func Test(t *testing.T) {
	err := errors.New("error")
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
	err := errors.New("error")
	ignored := Is(err)
	if ignored {
		t.Fatalf("unexpected ignored: got %t, want %t", ignored, false)
	}
}

func TestError(t *testing.T) {
	err := errors.New("error")
	err = Wrap(err)
	s := err.Error()
	expected := "error"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
}

func TestVerbose(t *testing.T) {
	err := errors.New("error")
	err = Wrap(err)
	var v errors.Verboser
	ok := errors.As(err, &v)
	if !ok {
		t.Fatal("not a Verbose")
	}
	buf := new(bytes.Buffer)
	v.ErrorVerbose(buf)
	s := buf.String()
	expected := "ignored\n"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
}

func TestUnwrap(t *testing.T) {
	err1 := errors.New("error")
	err2 := Wrap(err1)
	err2 = errors.Unwrap(err2)
	if err2 != err1 { //nolint:errorlint // We want to compare the error.
		t.Fatal("error not equal")
	}
}
