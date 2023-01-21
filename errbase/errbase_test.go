package errbase_test

import (
	"testing"

	. "github.com/pierrre/errors/errbase"
)

func Test(t *testing.T) {
	err := Newf("error %d", 1)
	s := err.Error()
	expected := "error 1"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
}
