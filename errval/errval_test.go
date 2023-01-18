package errval

import (
	"fmt"
	"testing"

	"github.com/pierrre/errors"
	"github.com/pierrre/errors/errverbose"
)

func Test(t *testing.T) {
	err := errors.New("error")
	err = Wrap(err, "foo", "bar")
	vals := Get(err)
	if len(vals) != 1 {
		t.Fatalf("unexpected length: got %d, want %d", len(vals), 1)
	}
	if vals["foo"] != "bar" {
		t.Fatalf("unexpected value: got %q, want %q", vals["foo"], "bar")
	}
}

func TestOverWrite(t *testing.T) {
	err := errors.New("error")
	err = Wrap(err, "test", 1)
	err = Wrap(err, "test", 2)
	vals := Get(err)
	if len(vals) != 1 {
		t.Fatalf("unexpected length: got %d, want %d", len(vals), 1)
	}
	if vals["test"] != 2 {
		t.Fatalf("unexpected value: got %v, want %d", vals["test"], 1)
	}
}

func TestNil(t *testing.T) {
	err := Wrap(nil, "foo", "bar")
	if err != nil {
		t.Fatal(err)
	}
}

func TestEmpty(t *testing.T) {
	err := errors.New("error")
	vals := Get(err)
	if len(vals) != 0 {
		t.Fatalf("values not empty: got %#v", vals)
	}
}

func TestError(t *testing.T) {
	err := errors.New("error")
	err = Wrap(err, "foo", "bar")
	s := err.Error()
	expected := "error"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
}

func TestVerbose(t *testing.T) {
	err := errors.New("error")
	err = Wrap(err, "foo", "bar")
	var v errverbose.Interface
	ok := errors.As(err, &v)
	if !ok {
		t.Fatal("not a Verbose")
	}
	s := v.ErrorVerbose()
	expected := "value foo = bar"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
}

func Example() {
	err := errors.New("error")
	err = Wrap(err, "foo", "bar")
	vals := Get(err)
	fmt.Println(vals["foo"])
	// Output: bar
}
