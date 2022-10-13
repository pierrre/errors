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
