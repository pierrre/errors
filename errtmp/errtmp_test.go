package errtmp

import (
	"bytes"
	"testing"

	"github.com/pierrre/errors"
)

func TestTrue(t *testing.T) {
	err := errors.New("error")
	err = Wrap(err, true)
	temporary := Is(err)
	if !temporary {
		t.Fatal("not temporary")
	}
}

func TestFalse(t *testing.T) {
	err := errors.New("error")
	err = Wrap(err, false)
	temporary := Is(err)
	if temporary {
		t.Fatal("temporary")
	}
}

func TestDefault(t *testing.T) {
	err := errors.New("error")
	temporary := Is(err)
	if !temporary {
		t.Fatal("not temporary")
	}
}

func TestNil(t *testing.T) {
	err := Wrap(nil, true)
	if err != nil {
		t.Fatal(err)
	}
}

func TestError(t *testing.T) {
	err := errors.New("error")
	err = Wrap(err, true)
	s := err.Error()
	expected := "error"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
}

func TestVerbose(t *testing.T) {
	err := errors.New("error")
	err = Wrap(err, true)
	var v errors.Verboser
	ok := errors.As(err, &v)
	if !ok {
		t.Fatal("not a Verbose")
	}
	buf := new(bytes.Buffer)
	v.ErrorVerbose(buf)
	s := buf.String()
	expected := "temporary = true\n"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
}

func TestUnwrap(t *testing.T) {
	err1 := errors.New("error")
	err2 := Wrap(err1, true)
	err2 = errors.Unwrap(err2)
	if err2 != err1 { //nolint:errorlint // We want to compare the error.
		t.Fatal("error not equal")
	}
}
