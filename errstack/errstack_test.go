package errstack_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/pierrre/errors"
	"github.com/pierrre/errors/errbase"
	. "github.com/pierrre/errors/errstack"
	"github.com/pierrre/errors/errverbose"
)

func Test(t *testing.T) {
	err := errbase.New("error")
	err = Ensure(err)
	err = Ensure(err)
	sfs := Frames(err)
	if len(sfs) != 1 {
		t.Fatalf("unexpected length: got %d, want %d", len(sfs), 1)
	}
	sf := sfs[0]
	if sf == nil {
		t.Fatal("no stack frames")
	}
	f, _ := sf.Next()
	expectedFunction := "github.com/pierrre/errors/errstack_test.Test"
	if f.Function != expectedFunction {
		t.Fatalf("unexpected function: got %q, want %q", f.Function, expectedFunction)
	}
}

func TestNil(t *testing.T) {
	err := Wrap(nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestError(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err)
	s := err.Error()
	expected := "error"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
}

func TestVerbose(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err)
	var v errverbose.Interface
	ok := errors.As(err, &v)
	if !ok {
		t.Fatal("not a Verbose")
	}
	s := v.ErrorVerbose()
	expectedRegexp := regexp.MustCompile(`^stack\n(\t.+ .+:\d+\n)+$`)
	if !expectedRegexp.MatchString(s) {
		t.Fatalf("unexpected verbose message:\ngot: %q\nwant match: %q", s, expectedRegexp)
	}
}

func TestStackFrames(t *testing.T) {
	err := errbase.New("error")
	err = Wrap(err)
	var sErr interface {
		StackFrames() []uintptr
	}
	ok := errors.As(err, &sErr)
	if !ok {
		t.Fatal("not a stack")
	}
	pcs := sErr.StackFrames()
	if len(pcs) == 0 {
		t.Fatal("no stack PCs")
	}
}

func Example() {
	err := errors.New("error")
	err = Wrap(err)
	fmt.Println(err)
	sfs := Frames(err)
	fmt.Println(len(sfs))
	// Output:
	// error
	// 2
}
