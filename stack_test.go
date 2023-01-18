package errors_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/pierrre/errors"
	"github.com/pierrre/errors/errbase"
	"github.com/pierrre/errors/errverbose"
)

func TestStack(t *testing.T) {
	err := errbase.New("error")
	err = errors.Stack(err)
	sfs := errors.StackFrames(err)
	if len(sfs) != 1 {
		t.Fatalf("unexpected length: got %d, want %d", len(sfs), 1)
	}
	sf := sfs[0]
	if sf == nil {
		t.Fatal("no stack frames")
	}
	f, _ := sf.Next()
	expectedFunction := "github.com/pierrre/errors_test.TestStack"
	if f.Function != expectedFunction {
		t.Fatalf("unexpected function: got %q, want %q", f.Function, expectedFunction)
	}
}

func TestStackNil(t *testing.T) {
	err := errors.Stack(nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStackError(t *testing.T) {
	err := errbase.New("error")
	err = errors.Stack(err)
	s := err.Error()
	expected := "error"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
}

func TestStackVerbose(t *testing.T) {
	err := errbase.New("error")
	err = errors.Stack(err)
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
	err = errors.Stack(err)
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

func ExampleStack() {
	err := errors.New("error")
	err = errors.Stack(err)
	fmt.Println(err)
	sfs := errors.StackFrames(err)
	fmt.Println(len(sfs))
	// Output:
	// error
	// 2
}
