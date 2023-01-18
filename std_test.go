package errors_test

import (
	"io/fs"
	"testing"

	. "github.com/pierrre/errors"
	"github.com/pierrre/errors/errbase"
)

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
