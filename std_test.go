package errors

import (
	"io"
	"testing"
)

func TestIs(t *testing.T) {
	err := io.EOF
	err = Wrap(err, "test")
	ok := Is(err, io.EOF)
	if !ok {
		t.Fatal("not ok")
	}
}

func TestJoin(t *testing.T) {
	err := Join(New("error 1"), New("error 2"))
	errUnwrap, ok := err.(interface { //nolint:errorlint // We want to check this error object.
		Unwrap() []error
	})
	if !ok {
		t.Fatal("not ok")
	}
	errs := errUnwrap.Unwrap()
	if len(errs) != 2 {
		t.Fatalf("unexpected length: got %d, want %d", len(errs), 2)
	}
}
