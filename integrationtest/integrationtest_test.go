package integrationtest

import (
	"regexp"
	"testing"

	"github.com/pierrre/errors"
	"github.com/pierrre/errors/errignore"
	"github.com/pierrre/errors/errstack"
	"github.com/pierrre/errors/errtag"
	"github.com/pierrre/errors/errtmp"
	"github.com/pierrre/errors/errval"
	"github.com/pierrre/errors/errverbose"
)

func Test(t *testing.T) {
	err := errors.New("error")
	err = errors.Wrap(err, "test")
	err = errignore.Wrap(err)
	err = errtmp.Wrap(err, true)
	err = errtag.Wrap(err, "a", "b")
	err = errval.Wrap(err, "c", "d")
	for _, tc := range []struct {
		name string
		f    func(*testing.T, error)
	}{
		{"Error", testError},
		{"Verbose", testVerbose},
		{"Stack", testStack},
		{"Ignore", testIgnore},
		{"Temporary", testTemporary},
		{"Tag", testTag},
		{"Value", testValue},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc.f(t, err)
		})
	}
}

func testError(t *testing.T, err error) {
	t.Helper()
	s := err.Error()
	expected := "test: error"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
}

func testVerbose(t *testing.T, err error) {
	t.Helper()
	s := errverbose.String(err)
	expected := regexp.MustCompile(`^test: error\nvalue c = d\ntag a = b\ntemporary = true\nignored\nstack\n(\t.+ .+:\d+\n)+\n$`)
	if !expected.MatchString(s) {
		t.Fatalf("unexpected verbose message:\ngot: %q\nwant match: %q", s, expected)
	}
}

func testStack(t *testing.T, err error) {
	t.Helper()
	sfs := errstack.Frames(err)
	if len(sfs) != 1 {
		t.Fatalf("unexpected stack frames length: got %d, want 1", len(sfs))
	}
	sf := sfs[0]
	f, _ := sf.Next()
	expectedFunction := "github.com/pierrre/errors/integrationtest.Test"
	if f.Function != expectedFunction {
		t.Fatalf("unexpected function: got %q, want %q", f.Function, expectedFunction)
	}
}

func testIgnore(t *testing.T, err error) {
	t.Helper()
	if !errignore.Is(err) {
		t.Fatal("not ignored")
	}
}

func testTemporary(t *testing.T, err error) {
	t.Helper()
	if !errtmp.Is(err) {
		t.Fatal("not temporary")
	}
}

func testTag(t *testing.T, err error) {
	t.Helper()
	tags := errtag.Get(err)
	if len(tags) != 1 {
		t.Fatalf("unexpected length: got %d, want %d", len(tags), 1)
	}
	if tags["a"] != "b" {
		t.Fatalf("unexpected tag: got %q, want %q", tags["a"], "b")
	}
}

func testValue(t *testing.T, err error) {
	t.Helper()
	values := errval.Get(err)
	if len(values) != 1 {
		t.Fatalf("unexpected length: got %d, want %d", len(values), 1)
	}
	if values["c"] != "d" {
		t.Fatalf("unexpected tag: got %v, want %q", values["c"], "d")
	}
}
