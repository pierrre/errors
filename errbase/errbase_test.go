package errbase_test

import (
	"testing"

	. "github.com/pierrre/errors/errbase"
)

func Test(t *testing.T) {
	err := New("error")
	s := err.Error()
	expected := "error"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
}
