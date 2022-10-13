package errors

import (
	"testing"
)

func TestBase(t *testing.T) {
	err := newBase("error")
	s := err.Error()
	expected := "error"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
}
