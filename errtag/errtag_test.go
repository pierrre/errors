package errtag

import (
	"fmt"
	"testing"

	"github.com/pierrre/errors"
	"github.com/pierrre/errors/errverbose"
)

func Test(t *testing.T) {
	err := errors.New("error")
	err = Wrap(err, "foo", "bar")
	tags := Get(err)
	if len(tags) != 1 {
		t.Fatalf("unexpected length: got %d, want %d", len(tags), 1)
	}
	if tags["foo"] != "bar" {
		t.Fatalf("unexpected tag: got %q, want %q", tags["foo"], "bar")
	}
}

func Example() {
	err := errors.New("error")
	err = Wrap(err, "foo", "bar")
	tags := Get(err)
	fmt.Println(tags["foo"])
	// Output: bar
}

func TestInt(t *testing.T) {
	err := errors.New("error")
	err = WrapInt(err, "foo", 123)
	tags := Get(err)
	if len(tags) != 1 {
		t.Fatalf("unexpected length: got %d, want %d", len(tags), 1)
	}
	if tags["foo"] != "123" {
		t.Fatalf("unexpected tag: got %q, want %q", tags["foo"], "123")
	}
}

func TestInt64(t *testing.T) {
	err := errors.New("error")
	err = WrapInt64(err, "foo", 123)
	tags := Get(err)
	if len(tags) != 1 {
		t.Fatalf("unexpected length: got %d, want %d", len(tags), 1)
	}
	if tags["foo"] != "123" {
		t.Fatalf("unexpected tag: got %q, want %q", tags["foo"], "123")
	}
}

func TestFloat64(t *testing.T) {
	err := errors.New("error")
	err = WrapFloat64(err, "foo", 12.3)
	tags := Get(err)
	if len(tags) != 1 {
		t.Fatalf("unexpected length: got %d, want %d", len(tags), 1)
	}
	if tags["foo"] != "12.3" {
		t.Fatalf("unexpected tag: got %q, want %q", tags["foo"], "12.3")
	}
}

func TestBool(t *testing.T) {
	err := errors.New("error")
	err = WrapBool(err, "foo", true)
	tags := Get(err)
	if len(tags) != 1 {
		t.Fatalf("unexpected length: got %d, want %d", len(tags), 1)
	}
	if tags["foo"] != "true" {
		t.Fatalf("unexpected tag: got %q, want %q", tags["foo"], "true")
	}
}

func TestOverWrite(t *testing.T) {
	err := errors.New("error")
	err = Wrap(err, "test", "1")
	err = Wrap(err, "test", "2")
	tags := Get(err)
	if len(tags) != 1 {
		t.Fatalf("unexpected length: got %d, want %d", len(tags), 1)
	}
	if tags["test"] != "2" {
		t.Fatalf("unexpected tag: got %q, want %q", tags["test"], "2")
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
	tags := Get(err)
	if len(tags) != 0 {
		t.Fatalf("tags not empty: got %#v", tags)
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
	expected := "tag foo = bar"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
}
